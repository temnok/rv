package rv

type CPU struct {
	x          [32]int32
	pc, nextPC int32
	csr        CSR
	mem        []byte

	trapped bool

	reserved        bool
	reservedAddress int32

	priv int32
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
	// ExceptionInstructionAddressMisaligned = 0
	ExceptionInstructionAccessFault = 1
	ExceptionIllegalIstruction      = 2
	ExceptionBreakpoint             = 3
	// ExceptionLoadAddressMisaligned        = 4
	ExceptionLoadAccessFault = 5
	// ExceptionStoreAMOAddressMisaligned = 6
	ExceptionStoreAMOAccessFault      = 7
	ExceptionEnvironmentCallFromUMode = 8
	ExceptionEnvironmentCallFromSMode = 9
	ExceptionEnvironmentCallFromMMode = 11
	// ExceptionInstructionPageFault         = 12
	// ExceptionLoadPageFault                = 13
	// ExceptionStoreAMOPageFault            = 15

	// PrivU = 0
	PrivS = 1
	PrivM = 3
)

func (cpu *CPU) init(ramSize int) {
	const xlen32bit = 0b_01

	*cpu = CPU{
		pc:   ramBaseAddr,
		mem:  make([]byte, ramSize),
		priv: PrivM,
		csr: CSR{
			misa: xlen32bit<<30 |
				1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
				1<<('u'-'a') | 1<<('s'-'a'),
		},
	}
}

func (cpu *CPU) step() {
	cpu.trapped = false

	if cpu.trapOnPendingInterrupts(); cpu.trapped {
		cpu.pc = cpu.nextPC
		return
	}

	if opcode := cpu.memFetch(cpu.pc); !cpu.trapped {

		if compressed := opcode&0b11 != 0b11; compressed {
			opcode = decompress(opcode)
			cpu.nextPC = cpu.pc + 2
		} else {
			cpu.nextPC = cpu.pc + 4
		}

		cpu.exec(opcode)
	}

	cpu.pc = cpu.nextPC
}

func (cpu *CPU) trapCause() int32 {
	switch cpu.priv {
	case PrivM:
		return cpu.csr.mcause
	case PrivS:
		return cpu.csr.scause
	default:
		return -1
	}
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapOnPendingInterrupts() {
	mi := cpu.csr.mip & cpu.csr.mie

	if cpu.trapped || mi == 0 {
		return
	}

	for i := int32(12); i > 0; i-- {
		if bit(mi, i) == 0 {
			continue
		}

		priv := int32(PrivM)
		if bit(cpu.csr.mideleg, i) != 0 {
			priv = PrivS
		}

		if (priv == cpu.priv && bit(cpu.csr.mstatus, priv) != 0) || priv > cpu.priv {
			cpu.trap(-1<<31 | i)
			return
		}
	}
}

func (cpu *CPU) trap(cause int32) {
	cpu.trapped = true

	isInterrupt := bit(cause, 31) != 0
	causeID := bits(cause, 0, 5)

	deleg := cpu.csr.medeleg
	if isInterrupt {
		deleg = cpu.csr.mideleg
	}

	effectivePriv := int32(PrivM)
	if cpu.priv < PrivM && bit(deleg, causeID) != 0 {
		effectivePriv = PrivS
	}

	var tvec int32

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.csr.mstatus, mstatusMIE)
		cpu.csr.mstatus &^= 0b_11<<mstatusMPP | 1<<mstatusMPIE | 1<<mstatusMIE
		cpu.csr.mstatus |= cpu.priv<<mstatusMPP | mie<<mstatusMPIE

		cpu.csr.mepc = cpu.pc
		cpu.csr.mcause = cause
		cpu.csr.mtval = 0
		tvec = cpu.csr.mtvec

	case PrivS:
		sie := bit(cpu.csr.mstatus, mstatusSIE)
		cpu.csr.mstatus &^= 1<<mstatusSPP | 1<<mstatusSPIE | 1<<mstatusSIE
		cpu.csr.mstatus |= cpu.priv<<mstatusSPP | sie<<mstatusSPIE

		cpu.csr.sepc = cpu.pc
		cpu.csr.scause = cause
		cpu.csr.stval = 0
		tvec = cpu.csr.stvec
	}

	cpu.nextPC = tvec &^ 0b_11
	if bit(tvec, 0) != 0 && isInterrupt {
		cpu.nextPC += causeID * 4
	}

	cpu.priv = effectivePriv
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) ret(priv int32) {
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
