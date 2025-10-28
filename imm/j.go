package imm

import "github.com/temnok/rv/bi"

func J(opcode int) int {
	a := bi.Ts(opcode, 21, 10)
	b := bi.T(opcode, 20)
	c := bi.Ts(opcode, 12, 8)
	d := bi.T(opcode, 31)

	return -d<<20 | c<<12 | b<<11 | a<<1
}
