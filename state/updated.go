package state

const (
	UpdateReservation = 1
)

type UpdatedState struct {
	Targets int

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

	Reservation int
}
