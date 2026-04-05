package imm

func CSS3(op int) int {
	a := op >> 7 & 7
	b := op >> 10 & 7

	return a<<6 | b<<3
}
