package imm

func CI(instr int) int {
	a := bits(instr, 2, 5)
	b := signBit(instr, 12)

	return b<<5 | a
}
