package imm

import "github.com/temnok/rv/bi"

func CB(instr int) int {
	a := bi.Ts(instr, 3, 2)
	b := bi.Ts(instr, 10, 2)
	c := bi.T(instr, 2)
	d := bi.Ts(instr, 5, 2)
	e := bi.T(instr, 12)

	return -e<<8 | d<<6 | c<<5 | b<<3 | a<<1
}
