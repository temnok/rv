package imm

import "github.com/temnok/rv/bi"

func CSS3(instr int) int {
	a := bi.Ts(instr, 10, 3)
	b := bi.Ts(instr, 7, 3)

	return b<<6 | a<<3
}
