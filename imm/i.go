package imm

func I(op int) int {
	a := op >> 31 & 1
	b := op >> 20 & 0x7FF

	return -a<<11 | b
}
