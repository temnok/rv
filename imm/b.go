package imm

func B(op int) int {
	a := op >> 31 & 1
	b := op >> 7 & 1
	c := op >> 25 & 0x3F
	d := op >> 8 & 0xF

	return -a<<12 | b<<11 | c<<5 | d<<1
}
