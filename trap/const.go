package trap

// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:mcause_exccode_enc_img
const (
	InstructionAddressMisaligned = 0
	InstructionAccessFault       = 1
	IllegalIstruction            = 2
	Breakpoint                   = 3
	LoadAddressMisaligned        = 4
	LoadAccessFault              = 5
	StoreAMOAddressMisaligned    = 6
	StoreAMOAccessFault          = 7
	EnvironmentCallFromUMode     = 8
	EnvironmentCallFromSMode     = 9
	EnvironmentCallFromMMode     = 11
	PageFault                    = 12
	InstructionPageFault         = 12
	LoadPageFault                = 13
	StoreAMOPageFault            = 15
	DoubleTrap                   = 16
	SoftwareCheck                = 18
	HardwareError                = 19
)
