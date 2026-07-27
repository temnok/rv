package encode

func B(imm, rs2, rs1, f3, op int) int {
	imm4_1 := imm >> 1 & 0xF
	imm10_5 := imm >> 5 & 0x3F
	imm11 := imm >> 11 & 1
	imm12 := imm >> 12 & 1

	return imm12<<31 | imm10_5<<25 | rs2<<20 | rs1<<15 | f3<<12 | imm4_1<<8 | imm11<<7 | op<<2 | 3
}
