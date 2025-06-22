package rv

type CPU struct {
	x          [32]int32
	pc, nextPC int32
	csr        CSR
	mem        []byte

	trapped bool
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

func (cpu *CPU) trap(cause int32) {
	cpu.csr.mepc = cpu.pc
	cpu.csr.mcause = cause
	cpu.csr.mtval = 0
	cpu.nextPC = cpu.csr.mtvec &^ 0b_11
	cpu.pc = cpu.nextPC

	cpu.trapped = true
}
