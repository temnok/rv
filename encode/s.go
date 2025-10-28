package encode

func S(imm, rs2, rs1, f3, op int) int {
	a := bits(imm, 5, 7)
	b := bits(imm, 0, 5)
	return a<<25 | rs2<<20 | rs1<<15 | f3<<12 | b<<7 | op<<2 | 3
}
