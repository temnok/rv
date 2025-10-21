package imm

func CI3(instr int) int {
	a := bits(instr, 5, 2)
	b := bit(instr, 12)
	c := bits(instr, 2, 3)

	return c<<6 | b<<5 | a<<3
}
