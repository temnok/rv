package imm

func CI2(op int) int {
	a := op >> 2 & 3
	b := op >> 12 & 1
	c := op >> 4 & 7

	return a<<6 | b<<5 | c<<2
}
