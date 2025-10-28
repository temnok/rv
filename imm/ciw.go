package imm

import "github.com/temnok/rv/bi"

func CIW(instr int) int {
	a := bi.T(instr, 6)
	b := bi.T(instr, 5)
	c := bi.Ts(instr, 11, 2)
	d := bi.Ts(instr, 7, 4)

	return d<<6 | c<<4 | b<<3 | a<<2
}
