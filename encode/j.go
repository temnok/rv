package encode

import "github.com/temnok/rv/bi"

func J(imm, rd, op int) int {
	a := bi.T(imm, 20)
	b := bi.Ts(imm, 12, 8)
	c := bi.T(imm, 11)
	d := bi.Ts(imm, 1, 10)
	return a<<31 | d<<21 | c<<20 | b<<12 | rd<<7 | op<<2 | 3
}
