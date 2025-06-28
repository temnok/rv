package rv

type (
	Xint  = int32
	Xuint = uint32
)

const (
	Xshift = 5
	Xbits  = 1 << Xshift
	Xbytes = Xbits / 8
)
