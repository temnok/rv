package imm

import "github.com/temnok/rv/bi"

func CI4(instr int) int {
	a := bi.T(instr, 6)
	b := bi.T(instr, 2)
	c := bi.T(instr, 5)
	d := bi.Ts(instr, 3, 2)
	e := bi.T(instr, 12)

	return -e<<9 | d<<7 | c<<6 | b<<5 | a<<4
}
