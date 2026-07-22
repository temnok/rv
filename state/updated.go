package state

const (
	UpdatePriv        = 1 << 0
	UpdateMstatus     = 1 << 1
	UpdateReservation = 1 << 2
	UpdateEpc         = 1 << 3
	UpdateXreg        = 1 << 4
	UpdateFreg        = 1 << 5
	UpdateCreg        = 1 << 6
)

type UpdatedState struct {
	Targets int

	Priv        int
	Reservation int

	TrapEnter  bool
	Mstatus    int
	TrapXcause int
	TrapXtval  int

	PC         int
	Xreg, Xval int

	Freg, Fval int
	Fflags     int

	Creg, Cval int
}
