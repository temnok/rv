package imm

import "github.com/temnok/rv/bi"

func CI2(instr int) int {
	a := bi.Ts(instr, 4, 3)
	b := bi.T(instr, 12)
	c := bi.Ts(instr, 2, 2)

	return c<<6 | b<<5 | a<<2
}
