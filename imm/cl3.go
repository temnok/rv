package imm

func CL3(op int) int {
	a := op >> 5 & 3
	b := op >> 10 & 7

	return a<<6 | b<<3
}
