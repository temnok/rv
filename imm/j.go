package imm

func J(op int) int {
	a := op >> 31 & 1
	b := op >> 12 & 255
	c := op >> 20 & 1
	d := op >> 21 & 1023

	return -a<<20 | b<<12 | c<<11 | d<<1
}
