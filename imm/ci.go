package imm

import "github.com/temnok/rv/bi"

func CI(instr int) int {
	a := bi.Ts(instr, 2, 5)
	b := bi.T(instr, 12)

	return -b<<5 | a
}
