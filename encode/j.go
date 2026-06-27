package encode

func J(imm, rd, op int) int {
	imm10_1 := imm >> 1 & 1023
	imm11 := imm >> 11 & 1
	imm19_12 := imm >> 12 & 255
	imm20 := imm >> 20 & 1

	return imm20<<31 | imm10_1<<21 | imm11<<20 | imm19_12<<12 | rd<<7 | op<<2 | 3
}
