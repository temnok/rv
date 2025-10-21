package imm

func CL3(instr int) int {
	a := bits(instr, 10, 3)
	b := bits(instr, 5, 2)

	return b<<6 | a<<3
}
