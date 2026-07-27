package imm

func CI4(op int) int {
	a := op >> 12 & 1
	b := op >> 3 & 3
	c := op >> 5 & 1
	d := op >> 2 & 1
	e := op >> 6 & 1

	return -a<<9 | b<<7 | c<<6 | d<<5 | e<<4
}
