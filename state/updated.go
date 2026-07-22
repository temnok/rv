package state

const (
	UpdateXreg        = 1 << 0
	UpdateFreg        = 1 << 1
	UpdateCreg        = 1 << 2
	UpdatePriv        = 1 << 3
	UpdateMstatus     = 1 << 4
	UpdateEpc         = 1 << 5
	UpdateCause       = 1 << 6
	UpdateTval        = 1 << 7
	UpdateFflags      = 1 << 8
	UpdateReservation = 1 << 9
)

type UpdatedState struct {
	Targets int

	Priv        int
	Reservation int

	TrapEnter bool
	Mstatus   int
	Cause     int
	Tval      int

	PC         int
	Xreg, Xval int

	Freg, Fval int
	Fflags     int

	Creg, Cval int
}
