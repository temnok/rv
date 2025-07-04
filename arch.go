package rv

type (
	Xint  = int64  // rv32: int32, rv64: int64
	Xuint = uint64 // rv32: uint32, rv64: uint64
)

const (
	Xshift = 6 // rv32: 5, rv64: 6

	Xlen   = 1 << Xshift
	Xbytes = Xlen / 8

	Xlen32 = Xlen == 32
	Xlen64 = Xlen == 64
)
