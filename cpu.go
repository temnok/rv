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
	Reg  [32]int
	CSR  CSR

	Reserved        bool
	ReservedAddress int

	ICache Cache
}

type CPUUpdatedState struct {
	Priv               int
	PC                 int
	RegIndex, RegValue int

	CSRAddr            *int
	CSRIndex, CSRValue int

	TrapState                       bool
	TrapXepc, TrapXcause, TrapXtval int

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
			CSR: CSR{
				Misa: xl<<(xlen-2) |
					1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
					1<<('u'-'a') | 1<<('s'-'a'),
			},
		},

		Updated: CPUUpdatedState{
			Priv:     PrivM,
			RegIndex: regIndex,
			RegValue: regValue,
		},
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

func (cpu *CPU) Step() {
	cpu.updateState()

	cpu.Updated.Priv = cpu.Priv
	cpu.Updated.Reserved = cpu.Reserved
	cpu.Updated.ReservedAddress = cpu.ReservedAddress

	if cpu.trapOnPendingInterrupts(); cpu.isTrapped() {
		return
	}

	var opcode int
	if cpu.memFetch(cpu.PC, &opcode); cpu.isTrapped() {
		return
	}

	cpu.exec(opcode)
}

func (cpu *CPU) isTrapped() bool {
	return cpu.Updated.TrapState
}

func (cpu *CPU) updateState() {
	cpu.updateTimers()
	cpu.clearPendingInterrupts()

	up := &cpu.Updated

	cpu.PC = up.PC
	cpu.Priv = up.Priv

	if up.RegIndex != 0 {
		cpu.Reg[up.RegIndex] = up.RegValue
	}

	if up.CSRIndex != 0 {
		*up.CSRAddr = up.CSRValue
	}

	if up.TrapState {
		if up.Priv == PrivM {
			cpu.CSR.Mepc = up.TrapXepc
			cpu.CSR.Mcause = up.TrapXcause
			cpu.CSR.Mtval = up.TrapXtval
		} else {
			cpu.CSR.Sepc = up.TrapXepc
			cpu.CSR.Scause = up.TrapXcause
			cpu.CSR.Stval = up.TrapXtval
		}
	}

	cpu.Reserved = up.Reserved
	if up.Reserved {
		cpu.ReservedAddress = up.ReservedAddress
	}

	cpu.ICache = cpu.Updated.ICache

	*up = CPUUpdatedState{}
}

func (cpu *CPU) updateTimers() {
	if cpu.CSR.Cycle = cpu.xint(cpu.CSR.Cycle + 1); cpu.CSR.Cycle == 0 {
		cpu.CSR.Cycleh++
	}

	if cpu.CSR.Cycle%10_000 == 0 {
		if cpu.CSR.Time = cpu.xint(cpu.CSR.Time + 1); cpu.CSR.Time == 0 {
			cpu.CSR.Timeh++
		}
	}
}

func (cpu *CPU) clearPendingInterrupts() {
	cpu.CSR.Mip &^= 1<<MipSEI | 1<<MipMTI | 1<<MipMSI
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

func (cpu *CPU) trap(cause int) {
	cpu.trapWithTval(cause, 0)
}

func (cpu *CPU) trapWithTval(cause, tval int) {
	if cpu.isTrapped() {
		panic("double trap")
	}

	mcauseI := cpu.XLen - 1
	isInterrupt := bit(cause, mcauseI) == 1
	causeID := bits(cause, 0, 5)

	deleg := cpu.CSR.Medeleg
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	}

	effectivePriv := PrivM
	if cpu.Priv < PrivM && bit(deleg, causeID) == 1 {
		effectivePriv = PrivS
	}

	var tvec int

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.CSR.Mstatus, MstatusMIE)
		cpu.Updated.CSRValue = cpu.CSR.Mstatus&^(3<<MstatusMPP|1<<MstatusMPIE|1<<MstatusMIE) |
			(cpu.Priv<<MstatusMPP | mie<<MstatusMPIE)

		tvec = cpu.CSR.Mtvec

	case PrivS:
		sie := bit(cpu.CSR.Mstatus, MstatusSIE)
		cpu.Updated.CSRValue = cpu.CSR.Mstatus&^(1<<MstatusSPP|1<<MstatusSPIE|1<<MstatusSIE) |
			(cpu.Priv<<MstatusSPP | sie<<MstatusSPIE)

		tvec = cpu.CSR.Stvec
	}

	cpu.Updated.PC = tvec &^ 3
	if bit(tvec, 0) == 1 && isInterrupt {
		cpu.Updated.PC += causeID * 4
	}

	cpu.Updated.TrapState = true
	cpu.Updated.Priv = effectivePriv
	cpu.Updated.CSRAddr = &cpu.CSR.Mstatus
	cpu.Updated.CSRIndex = Mstatus
	cpu.Updated.TrapXepc = cpu.PC
	cpu.Updated.TrapXcause = cause
	cpu.Updated.TrapXtval = tval
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) xret(priv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trap := cpu.Priv == PrivS && bit(cpu.CSR.Mstatus, MstatusTSR) == 1

	if trap || priv > cpu.Priv {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.CSRAddr = &cpu.CSR.Mstatus
	cpu.Updated.CSRIndex = Mstatus

	switch priv {
	case PrivM:
		cpu.Updated.PC = cpu.CSR.Mepc
		cpu.Updated.Priv = bits(cpu.CSR.Mstatus, MstatusMPP, 2)

		mie := bit(cpu.CSR.Mstatus, MstatusMPIE)
		cpu.Updated.CSRValue = cpu.CSR.Mstatus&^(3<<MstatusMPP) |
			(1<<MstatusMPIE | mie<<MstatusMIE)

	case PrivS:
		cpu.Updated.PC = cpu.CSR.Sepc
		cpu.Updated.Priv = bits(cpu.CSR.Mstatus, MstatusSPP, 1)

		sie := bit(cpu.CSR.Mstatus, MstatusSPIE)
		cpu.Updated.CSRValue = cpu.CSR.Mstatus&^(1<<MstatusSPP) |
			(1<<MstatusSPIE | sie<<MstatusSIE)
	}

	if cpu.Priv != PrivM {
		cpu.Updated.CSRValue &^= 1 << MstatusMPRV
	}
}
