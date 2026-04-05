package imm

func CL(op int) int {
	a := op >> 5 & 1
	b := op >> 10 & 7
	c := op >> 6 & 1

	return a<<6 | b<<3 | c<<2
}
