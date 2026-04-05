package imm

func U(op int) int {
	a := op >> 31 & 1
	b := op >> 12 & (1<<19 - 1)

	return -a<<31 | b<<12
}
