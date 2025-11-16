package state

type UpdatedState struct {
	TrapEnter   bool
	TrapExit    bool
	TrapMstatus int
	TrapPC      int
	TrapPriv    int
	TrapXcause  int
	TrapXepc    int
	TrapXtval   int

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
