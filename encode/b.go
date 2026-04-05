package encode

func B(imm, rs2, rs1, f3, op int) int {
	a := imm >> 12 & 1
	b := imm >> 11 & 1
	c := imm >> 5 & 63
	d := imm >> 1 & 15
	return a<<31 | c<<25 | rs2<<20 | rs1<<15 | f3<<12 | d<<8 | b<<7 | op<<2 | 3
}
