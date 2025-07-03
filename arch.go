package rv

type (
	Xint  = int32
	Xuint = uint32
)

const (
	Xshift = 5
	Xlen   = 1 << Xshift
	Xbytes = Xlen / 8

	XlenIs32 = Xlen == 32
	XlenIs64 = Xlen == 64
)
