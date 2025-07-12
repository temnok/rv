package rv

type CPU struct {
	Xlen int

	x          [32]int
	PC, nextPC int
	csr        CSR
	priv       int

	isTrapped bool

	reserved        bool
	reservedAddress int

	tlb TLB

	bus Bus
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
		Xlen: xlen,

		priv: PrivM,
		csr: CSR{
			misa: int(misa) |
				1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
				1<<('u'-'a') | 1<<('s'-'a'),
		},

		bus: bus,
	}

	cpu.PC = cpu.Xint(startAddr)
	cpu.csr.mstatus = cpu.Xint(xl<<mstatusSXL | xl<<mstatusUXL)

	for i, x := range regs {
		cpu.x[i] = cpu.Xint(x)
	}
}

func (cpu *CPU) Xlen64() bool {
	return cpu.Xlen == 64
}

func (cpu *CPU) Xint(val int) int {
	if cpu.Xlen64() {
		return val
	}

	return int(int32(val))
}

func (cpu *CPU) Xuint(val int) uint {
	if cpu.Xlen64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *CPU) Step() bool {
	//return cpu.debugStep()
	cpu.step()
	return true
}

func (cpu *CPU) step() int {
	cpu.isTrapped = false

	cpu.updateTimers()

	cpu.csr.mip &^= 1<<mipSEI | 1<<mipMTI | 1<<mipMSI
	cpu.bus.notifyInterrupts()
	if cpu.trapOnPendingInterrupts(); cpu.isTrapped {
		return 0
	}

	var opcode int
	if cpu.memFetch(cpu.PC, &opcode); cpu.isTrapped {
		return 0
	}

	cpu.nextPC = cpu.Xint(cpu.PC + 4)
	origOpcode := opcode
	if cpu.decompress(&opcode); cpu.isTrapped {
		return opcode
	}

	cpu.exec(opcode)

	return origOpcode
}

func (cpu *CPU) updateTimers() {
	if cpu.csr.cycle = cpu.Xint(cpu.csr.cycle + 1); cpu.csr.cycle == 0 {
		cpu.csr.cycleh++
	}

	if cpu.csr.cycle%10_000 == 0 {
		if cpu.csr.mtime = cpu.Xint(cpu.csr.mtime + 1); cpu.csr.mtime == 0 {
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
		if bit(cpu.csr.mideleg, i) != 0 {
			priv = PrivS
		}

		mcauseI := cpu.Xlen - 1
		if (priv == cpu.priv && bit(cpu.csr.mstatus, priv) != 0) || priv > cpu.priv {
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

	mcauseI := cpu.Xlen - 1
	isInterrupt := bit(cause, mcauseI) != 0
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
		cpu.csr.mstatus &^= 0b_11<<mstatusMPP | 1<<mstatusMPIE | 1<<mstatusMIE
		cpu.csr.mstatus |= cpu.priv<<mstatusMPP | mie<<mstatusMPIE

		cpu.csr.mepc = cpu.PC
		cpu.csr.mcause = cause
		cpu.csr.mtval = tval
		tvec = cpu.csr.mtvec

	case PrivS:
		sie := bit(cpu.csr.mstatus, mstatusSIE)
		cpu.csr.mstatus &^= 1<<mstatusSPP | 1<<mstatusSPIE | 1<<mstatusSIE
		cpu.csr.mstatus |= cpu.priv<<mstatusSPP | sie<<mstatusSPIE

		cpu.csr.sepc = cpu.PC
		cpu.csr.scause = cause
		cpu.csr.stval = tval
		tvec = cpu.csr.stvec
	}

	cpu.PC = tvec &^ 0b_11
	if bit(tvec, 0) != 0 && isInterrupt {
		cpu.PC += causeID * 4
	}

	cpu.priv = effectivePriv
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) ret(priv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapSRET := cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTSR) != 0

	if priv > cpu.priv || trapSRET {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch priv {
	case PrivM:
		cpu.nextPC = cpu.csr.mepc
		cpu.priv = bits(cpu.csr.mstatus, mstatusMPP, 2)

		mie := bit(cpu.csr.mstatus, mstatusMPIE)
		cpu.csr.mstatus |= 1<<mstatusMPIE | mie<<mstatusMIE
		cpu.csr.mstatus &^= 0b_11 << mstatusMPP

	case PrivS:
		cpu.nextPC = cpu.csr.sepc
		cpu.priv = bits(cpu.csr.mstatus, mstatusSPP, 1)

		sie := bit(cpu.csr.mstatus, mstatusSPIE)
		cpu.csr.mstatus |= 1<<mstatusSPIE | sie<<mstatusSIE
		cpu.csr.mstatus &^= 1 << mstatusSPP
	}

	if cpu.priv != PrivM {
		cpu.csr.mstatus &^= 1 << mstatusMPRV
	}
}
