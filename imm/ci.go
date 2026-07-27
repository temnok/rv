package imm

func CI(op int) int {
	a := op >> 12 & 1
	b := op >> 2 & 0x1F

	return -a<<5 | b
}
