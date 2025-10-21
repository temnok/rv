package imm

func CIW(instr int) int {
	a := bit(instr, 6)
	b := bit(instr, 5)
	c := bits(instr, 11, 2)
	d := bits(instr, 7, 4)

	return d<<6 | c<<4 | b<<3 | a<<2
}
