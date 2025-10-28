package imm

import "github.com/temnok/rv/bi"

func CL3(instr int) int {
	a := bi.Ts(instr, 10, 3)
	b := bi.Ts(instr, 5, 2)

	return b<<6 | a<<3
}
