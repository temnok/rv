package encode

func I(imm, rs1, f3, rd, op int) int {
	imm11_0 := imm & 0xFFF

	return imm11_0<<20 | rs1<<15 | f3<<12 | rd<<7 | op<<2 | 3
}
