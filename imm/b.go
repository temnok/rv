package imm

func B(opcode int) int {
	a := bits(opcode, 8, 4)
	b := bits(opcode, 25, 6)
	c := bit(opcode, 7)
	d := signBit(opcode, 31)

	return d<<12 | c<<11 | b<<5 | a<<1
}
