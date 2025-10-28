package imm

import "github.com/temnok/rv/bi"

func CI3(instr int) int {
	a := bi.Ts(instr, 5, 2)
	b := bi.T(instr, 12)
	c := bi.Ts(instr, 2, 3)

	return c<<6 | b<<5 | a<<3
}
