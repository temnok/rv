package imm

func CI3(op int) int {
	a := op >> 2 & 7
	b := op >> 12 & 1
	c := op >> 5 & 3

	return a<<6 | b<<5 | c<<3
}
