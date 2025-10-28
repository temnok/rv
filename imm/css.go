package imm

import "github.com/temnok/rv/bi"

func CSS(instr int) int {
	a := bi.Ts(instr, 9, 4)
	b := bi.Ts(instr, 7, 2)

	return b<<6 | a<<2
}
