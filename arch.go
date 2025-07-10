package rv

type (
	Xint  = int64
	Xuint = uint64
)

const (
	Xlen = 64

	Xbytes = Xlen / 8
	Xlen32 = Xlen == 32
	Xlen64 = Xlen == 64

	xshift = 64 - Xlen
)

func FixInt(val Xint) Xint {
	return val << xshift >> xshift
}

func FixUint(val Xuint) Xuint {
	return val << xshift >> xshift
}
