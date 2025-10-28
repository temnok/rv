package imm

import "github.com/temnok/rv/bi"

func CJ(instr int) int {
	a := bi.Ts(instr, 3, 3)
	b := bi.T(instr, 11)
	c := bi.T(instr, 2)
	d := bi.T(instr, 7)
	e := bi.T(instr, 6)
	f := bi.Ts(instr, 9, 2)
	g := bi.T(instr, 8)
	h := bi.T(instr, 12)

	return -h<<11 | g<<10 | f<<8 | e<<7 | d<<6 | c<<5 | b<<4 | a<<1
}
