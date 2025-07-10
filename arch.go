package rv

const (
	Xlen = 64

	Xbytes = Xlen / 8
	Xlen32 = Xlen == 32
	Xlen64 = Xlen == 64

	xshift = 64 - Xlen
)

func Xint(val int) int {
	return val << xshift >> xshift
}

func Xuint(val int) uint {
	return uint(val) << xshift >> xshift
}
