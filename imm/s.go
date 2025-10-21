package imm

func S(opcode int) int {
	a := bits(opcode, 7, 5)
	b := int(int32(opcode)) >> 25

	return b<<5 | a
}
