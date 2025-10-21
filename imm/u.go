package imm

func U(opcode int) int {
	return int(int32(bits(opcode, 12, 20) << 12))
}
