package rv

type (
	Xint  = int64
	Xuint = uint64
)

const (
	Xshift = 6
	Xlen   = 1 << Xshift
	Xbytes = Xlen / 8

	XlenIs32 = Xlen == 32
	XlenIs64 = Xlen == 64
)
