package imm

func S(op int) int {
	a := op >> 31 & 1
	b := op >> 25 & 63
	c := op >> 7 & 31

	return -a<<11 | b<<5 | c
}
