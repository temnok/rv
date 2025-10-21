package imm

func CL(instr int) int {
	a := bit(instr, 6)
	b := bits(instr, 10, 3)
	c := bit(instr, 5)

	return c<<6 | b<<3 | a<<2
}
