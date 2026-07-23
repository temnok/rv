package state

const (
	UpdatePC          = 1 << 0
	UpdatePriv        = 1 << 1
	UpdateReservation = 1 << 2

	UpdateXreg = 1 << 3
	UpdateFreg = 1 << 4
	UpdateCreg = 1 << 5

	UpdateFcsr    = 1 << 6
	UpdateMstatus = 1 << 7
	UpdateEpc     = 1 << 8
	UpdateCause   = 1 << 9
	UpdateTval    = 1 << 10
)

type Update struct {
	Targets int

	PC          int
	Priv        int
	Reservation int

	Xreg, Xval int
	Freg, Fval int
	Creg, Cval int

	Fcsr    int
	Mstatus int
	Epc     int
	Cause   int
	Tval    int
}
