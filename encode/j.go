package encode

func J(imm, rd, op int) int {
	a := imm >> 20 & 1
	b := imm >> 12 & 255
	c := imm >> 11 & 1
	d := imm >> 1 & 1023
	return a<<31 | d<<21 | c<<20 | b<<12 | rd<<7 | op<<2 | 3
}
