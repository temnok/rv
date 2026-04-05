package imm

func I(op int) int {
	a := op >> 31 & 1
	b := op >> 20 & 2047

	return -a<<11 | b
}
