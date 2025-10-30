package trap

const (
	InstructionAccessFault    = 1
	IllegalIstruction         = 2
	Breakpoint                = 3
	LoadAddressMisaligned     = 4
	LoadAccessFault           = 5
	StoreAMOAddressMisaligned = 6
	StoreAMOAccessFault       = 7
	EnvironmentCallFromUMode  = 8
	EnvironmentCallFromSMode  = 9
	EnvironmentCallFromMMode  = 11
	PageFault                 = 12
)
