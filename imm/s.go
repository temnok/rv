package imm

import "github.com/temnok/rv/bi"

func S(opcode int) int {
	a := bi.Ts(opcode, 7, 5)
	b := int(int32(opcode)) >> 25

	return b<<5 | a
}
