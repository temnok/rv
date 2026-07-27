package imm

func CIW(op int) int {
	a := op >> 7 & 0xF
	b := op >> 11 & 3
	c := op >> 5 & 1
	d := op >> 6 & 1

	return a<<6 | b<<4 | c<<3 | d<<2
}
