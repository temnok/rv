package imm

func CI2(instr int) int {
	a := bits(instr, 4, 3)
	b := bit(instr, 12)
	c := bits(instr, 2, 2)

	return c<<6 | b<<5 | a<<2
}
