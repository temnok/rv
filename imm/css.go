package imm

func CSS(op int) int {
	a := op >> 7 & 3
	b := op >> 9 & 0xF

	return a<<6 | b<<2
}
