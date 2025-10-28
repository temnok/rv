package imm

import "github.com/temnok/rv/bi"

func B(opcode int) int {
	a := bi.Ts(opcode, 8, 4)
	b := bi.Ts(opcode, 25, 6)
	c := bi.T(opcode, 7)
	d := bi.T(opcode, 31)

	return -d<<12 | c<<11 | b<<5 | a<<1
}
