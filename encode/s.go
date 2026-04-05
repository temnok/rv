package encode

func S(imm, rs2, rs1, f3, op int) int {
	a := imm >> 5 & 127
	b := imm & 31
	return a<<25 | rs2<<20 | rs1<<15 | f3<<12 | b<<7 | op<<2 | 3
}
