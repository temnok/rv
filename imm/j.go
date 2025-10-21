package imm

func J(opcode int) int {
	a := bits(opcode, 21, 10)
	b := bit(opcode, 20)
	c := bits(opcode, 12, 8)
	d := signBit(opcode, 31)

	return d<<20 | c<<12 | b<<11 | a<<1
}
