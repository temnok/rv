package imm

func I(opcode int) int {
	return int(int32(opcode) >> 20)
}
