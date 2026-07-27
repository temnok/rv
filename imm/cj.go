package imm

func CJ(op int) int {
	a := op >> 12 & 1
	b := op >> 8 & 1
	c := op >> 9 & 3
	d := op >> 6 & 1
	e := op >> 7 & 1
	f := op >> 2 & 1
	g := op >> 11 & 1
	h := op >> 3 & 7

	return -a<<11 | b<<10 | c<<8 | d<<7 | e<<6 | f<<5 | g<<4 | h<<1
}
