package rv

type CPU struct {
	xlen int

	reg  [32]int
	pc   int
	csr  CSR
	priv int

	isTrapped bool

	reserved        bool
	reservedAddress int

	tlb TLB

	bus Bus

	updated struct {
		pc       int
		priv     int
		regIndex int
		regValue int
	}
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
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

func (cpu *CPU) Init(xlen int, bus Bus, startAddr int, regs []int) {
	xl := xlen / 32
	misaMXL := xlen - 2
	misa := uint(xl << misaMXL)

	*cpu = CPU{
		xlen: xlen,

		csr: CSR{
			misa: int(misa) |
				1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
				1<<('u'-'a') | 1<<('s'-'a'),
		},

		bus: bus,
	}

	cpu.updated.priv = PrivM
	cpu.updated.pc = cpu.xint(startAddr)

	cpu.csr.mstatus = cpu.xint(xl<<mstatusSXL | xl<<mstatusUXL)

	for i, x := range regs {
		cpu.reg[i] = cpu.xint(x)
	}
}

func (cpu *CPU) xlen64() bool {
	return cpu.xlen == 64
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

func (cpu *CPU) step() bool {
	cpu.updateState()

	//return cpu.debugStep()
	cpu.innerStep()

	return true
}

func (cpu *CPU) innerStep() int {
	cpu.isTrapped = false

	cpu.updateTimers()

	cpu.csr.mip &^= 1<<mipSEI | 1<<mipMTI | 1<<mipMSI
	cpu.bus.notifyInterrupts()
	if cpu.trapOnPendingInterrupts(); cpu.isTrapped {
		return 0
	}

	var opcode int
	if cpu.memFetch(cpu.pc, &opcode); cpu.isTrapped {
		return 0
	}

	cpu.updated.pc = cpu.xint(cpu.pc + 4)

	origOpcode := opcode
	if cpu.decompress(&opcode); cpu.isTrapped {
		return opcode
	}

	cpu.exec(opcode)

	return origOpcode
}

func (cpu *CPU) updateState() {
	up := &cpu.updated

	cpu.pc = up.pc
	cpu.priv = up.priv

	if up.regIndex != 0 {
		cpu.reg[up.regIndex] = up.regValue
		up.regIndex = 0
	}
}

func (cpu *CPU) updateTimers() {
	if cpu.csr.cycle = cpu.xint(cpu.csr.cycle + 1); cpu.csr.cycle == 0 {
		cpu.csr.cycleh++
	}

	if cpu.csr.cycle%10_000 == 0 {
		if cpu.csr.mtime = cpu.xint(cpu.csr.mtime + 1); cpu.csr.mtime == 0 {
			cpu.csr.mtimeh++
		}
	}
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapOnPendingInterrupts() {
	mi := cpu.csr.mip & cpu.csr.mie

	if mi == 0 {
		return
	}

	for i := 12; i > 0; i-- {
		if bit(mi, i) == 0 {
			continue
		}

		priv := PrivM
		if bit(cpu.csr.mideleg, i) == 1 {
			priv = PrivS
		}

		mcauseI := cpu.xlen - 1
		if (priv == cpu.priv && bit(cpu.csr.mstatus, priv) == 1) || priv > cpu.priv {
			cpu.trap(-1<<mcauseI | i)

			return
		}
	}
}

func (cpu *CPU) trap(cause int) {
	cpu.trapWithTval(cause, 0)
}

func (cpu *CPU) trapWithTval(cause, tval int) {
	if cpu.isTrapped {
		panic("double trap")
	}

	cpu.isTrapped = true

	mcauseI := cpu.xlen - 1
	isInterrupt := bit(cause, mcauseI) == 1
	causeID := bits(cause, 0, 5)

	deleg := cpu.csr.medeleg
	if isInterrupt {
		deleg = cpu.csr.mideleg
	}

	effectivePriv := PrivM
	if cpu.priv < PrivM && bit(deleg, causeID) == 1 {
		effectivePriv = PrivS
	}

	var tvec int

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.csr.mstatus, mstatusMIE)
		cpu.csr.mstatus &^= 3<<mstatusMPP | 1<<mstatusMPIE | 1<<mstatusMIE
		cpu.csr.mstatus |= cpu.priv<<mstatusMPP | mie<<mstatusMPIE

		cpu.csr.mepc = cpu.pc
		cpu.csr.mcause = cause
		cpu.csr.mtval = tval
		tvec = cpu.csr.mtvec

	case PrivS:
		sie := bit(cpu.csr.mstatus, mstatusSIE)
		cpu.csr.mstatus &^= 1<<mstatusSPP | 1<<mstatusSPIE | 1<<mstatusSIE
		cpu.csr.mstatus |= cpu.priv<<mstatusSPP | sie<<mstatusSPIE

		cpu.csr.sepc = cpu.pc
		cpu.csr.scause = cause
		cpu.csr.stval = tval
		tvec = cpu.csr.stvec
	}

	cpu.updated.pc = tvec &^ 3
	if bit(tvec, 0) == 1 && isInterrupt {
		cpu.updated.pc += causeID * 4
	}

	cpu.updated.priv = effectivePriv
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) ret(priv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapSRET := cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTSR) == 1

	if priv > cpu.priv || trapSRET {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch priv {
	case PrivM:
		cpu.updated.pc = cpu.csr.mepc
		cpu.updated.priv = bits(cpu.csr.mstatus, mstatusMPP, 2)

		mie := bit(cpu.csr.mstatus, mstatusMPIE)
		cpu.csr.mstatus |= 1<<mstatusMPIE | mie<<mstatusMIE
		cpu.csr.mstatus &^= 3 << mstatusMPP

	case PrivS:
		cpu.updated.pc = cpu.csr.sepc
		cpu.updated.priv = bits(cpu.csr.mstatus, mstatusSPP, 1)

		sie := bit(cpu.csr.mstatus, mstatusSPIE)
		cpu.csr.mstatus |= 1<<mstatusSPIE | sie<<mstatusSIE
		cpu.csr.mstatus &^= 1 << mstatusSPP
	}

	if cpu.priv != PrivM {
		cpu.csr.mstatus &^= 1 << mstatusMPRV
	}
}
