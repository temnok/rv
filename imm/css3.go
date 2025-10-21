package imm

func CSS3(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 7, 3)

	return b<<6 | a<<3
}
