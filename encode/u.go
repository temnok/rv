package encode

func U(imm, rd, op int) int {
	return imm<<12 | rd<<7 | op<<2 | 3
}
