package state

type Updated struct {
	TrapEnter, TrapExit   bool
	TrapPriv, TrapPC      int
	TrapMstatus, TrapXepc int
	TrapXcause, TrapXtval int

	PC         int
	XReg, XVal int

	FReg, FVal int
	Fflags     int

	CReg, CVal int
	CRegPtr    *int

	Reserved     bool
	ReservedAddr int

	ICache Cache
}
