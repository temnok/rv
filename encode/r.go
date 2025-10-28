package encode

func R(f7, rs2, rs1, f3, rd, op int) int {
	return f7<<25 | rs2<<20 | rs1<<15 | f3<<12 | rd<<7 | op<<2 | 3
}
