package rv

type CPU struct {
	x          [32]int32
	pc, nextPC int32
	csr        CSR
	mem        []byte

	trapped bool

	reservationValid   bool
	reservationAddress int32
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcauses
const (
	ExceptionInstructionAddressMisaligned = 0
	ExceptionInstructionAccessFault       = 1
	ExceptionIllegalIstruction            = 2
	ExceptionBreakpoint                   = 3
	ExceptionLoadAddressMisaligned        = 4
	ExceptionLoadAccessFault              = 5
	ExceptionStoreAMOAddressMisaligned    = 6
	ExceptionStoreAMOAccessFault          = 7
	ExceptionEnvironmentCallFromUMode     = 8
	ExceptionEnvironmentCallFromSMode     = 9
	ExceptionEnvironmentCallFromMMode     = 11
	ExceptionInstructionPageFault         = 12
	ExceptionLoadPageFault                = 13
	ExceptionStoreAMOPageFault            = 15
)

func (cpu *CPU) step() {
	cpu.trapped = false

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

func (cpu *CPU) trap(cause int32) {
	cpu.trapped = true

	cpu.csr.mepc = cpu.pc
	cpu.csr.mcause = cause
	cpu.csr.mtval = 0

	cpu.nextPC = cpu.csr.mtvec &^ 0b_11
}
