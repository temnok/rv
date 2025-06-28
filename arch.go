package rv

type (
	Xint  = int64
	Xuint = uint64
)

const (
	Xshift = 6
	Xbits  = 1 << Xshift
	Xbytes = Xbits / 8
)
