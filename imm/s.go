package imm

func S(op int) int {
	a := op >> 31 & 1
	b := op >> 25 & 0x3F
	c := op >> 7 & 0x1F

	return -a<<11 | b<<5 | c
}
