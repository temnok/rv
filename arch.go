package rv

const (
	Xlen = 64

	Xbytes = Xlen / 8
	Xlen32 = Xlen == 32
	Xlen64 = Xlen == 64
)

func Xint(val int) int {
	if Xlen64 {
		return val
	}

	return int(int32(val))
}

func Xuint(val int) uint {
	if Xlen64 {
		return uint(val)
	}

	return uint(uint32(val))
}
