package encode

func J(imm, rd, op int) int {
	a := bit(imm, 20)
	b := bits(imm, 12, 8)
	c := bit(imm, 11)
	d := bits(imm, 1, 10)
	return a<<31 | d<<21 | c<<20 | b<<12 | rd<<7 | op<<2 | 3
}
