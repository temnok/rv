package rv

type CPU struct {
	XLen int
	Bus  Bus
	TLB  TLB

	CPUState
	Updated CPUUpdatedState
}

type CPUState struct {
	Priv int
	PC   int

	Reg, FReg [32]int

	CSR CSR

	Reserved        bool
	ReservedAddress int

	ICache Cache
}

type ICache struct {
	VirtAddr, PhysAddr, Value int
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
	PageSize = 1 << 12

	ExceptionInstructionAccessFault    = 1
	ExceptionIllegalIstruction         = 2
	ExceptionBreakpoint                = 3
	ExceptionLoadAddressMisaligned     = 4
	ExceptionLoadAccessFault           = 5
	ExceptionStoreAMOAddressMisaligned = 6
	ExceptionStoreAMOAccessFault       = 7
	ExceptionEnvironmentCallFromUMode  = 8
	ExceptionEnvironmentCallFromSMode  = 9
	ExceptionEnvironmentCallFromMMode  = 11
	ExceptionPageFault                 = 12

	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessExecute = 0
	AccessRead    = 1
	AccessWrite   = 3
)

func (cpu *CPU) Init(xlen int, bus Bus, startAddr, regIndex, regValue int) {
	xl := xlen / 32

	*cpu = CPU{
		XLen: xlen,
		Bus:  bus,

		CPUState: CPUState{
			Priv: PrivM,

			CSR: CSR{
				Misa: xl<<(xlen-2) |
					1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
					1<<('f'-'a') | (xlen/64)<<('d'-'a') |
					1<<('u'-'a') | 1<<('s'-'a'),
			},
		},

		Updated: CPUUpdatedState{
			RegIndex: -1,
			CSRIndex: -1,
		},
	}

	if regIndex > 0 {
		cpu.Updated.RegIndex = regIndex
		cpu.Updated.RegValue = regValue
	}

	cpu.CSR.Mstatus = cpu.xint(xl<<MstatusSXL | xl<<MstatusUXL)
	cpu.Updated.PC = cpu.xint(startAddr)
}

func (cpu *CPU) xlen64() bool {
	return cpu.XLen == 64
}

func (cpu *CPU) xint(val int) int {
	if cpu.xlen64() {
		return val
	}

	return int(int32(val))
}

func (cpu *CPU) xuint(val int) uint {
	if cpu.xlen64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *CPU) Step() bool {
	//return cpu.debugStep()

	cpu.innerStep()
	return true
}

func (cpu *CPU) innerStep() int {
	cpu.updateState()

	if cpu.trapOnPendingInterrupts(); cpu.isTrapped() {
		return 0
	}

	var opcode int
	if cpu.memFetch(cpu.PC, &opcode); cpu.isTrapped() {
		return 0
	}

	cpu.exec(opcode)

	return opcode
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapOnPendingInterrupts() {
	cpu.Bus.NotifyInterrupts()

	mi := cpu.CSR.Mip & cpu.CSR.Mie

	if mi == 0 {
		return
	}

	for i := 12; i > 0; i-- {
		if bit(mi, i) == 0 {
			continue
		}

		priv := PrivM
		if bit(cpu.CSR.Mideleg, i) == 1 {
			priv = PrivS
		}

		mcauseI := cpu.XLen - 1
		if (priv == cpu.Priv && bit(cpu.CSR.Mstatus, priv) == 1) || priv > cpu.Priv {
			cpu.trap(-1<<mcauseI | i)

			return
		}
	}
}
