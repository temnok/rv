package imm

import "github.com/temnok/rv/bi"

func U(opcode int) int {
	return int(int32(bi.Ts(opcode, 12, 20) << 12))
}
