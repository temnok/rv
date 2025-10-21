package imm

func CSS(instr int) int {
	a := bits(instr, 9, 4)
	b := bits(instr, 7, 2)

	return b<<6 | a<<2
}
