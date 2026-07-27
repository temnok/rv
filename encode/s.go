package encode

func S(imm, rs2, rs1, f3, op int) int {
	imm4_0 := imm & 0x1F
	imm11_5 := imm >> 5 & 0x7F

	return imm11_5<<25 | rs2<<20 | rs1<<15 | f3<<12 | imm4_0<<7 | op<<2 | 3
}
