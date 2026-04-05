package imm

func CB(op int) int {
	a := op >> 12 & 1
	b := op >> 5 & 3
	c := op >> 2 & 1
	d := op >> 10 & 3
	e := op >> 3 & 3

	return -a<<8 | b<<6 | c<<5 | d<<3 | e<<1
}
