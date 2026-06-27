package encode

func U(imm, rd, op int) int {
	imm31_12 := imm & 0xFFFFF

	return imm31_12<<12 | rd<<7 | op<<2 | 3
}
