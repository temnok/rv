package imm

import "github.com/temnok/rv/bi"

func CL(instr int) int {
	a := bi.T(instr, 6)
	b := bi.Ts(instr, 10, 3)
	c := bi.T(instr, 5)

	return c<<6 | b<<3 | a<<2
}
