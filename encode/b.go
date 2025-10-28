package encode

import "github.com/temnok/rv/bi"

func B(imm, rs2, rs1, f3, op int) int {
	a := bi.T(imm, 12)
	b := bi.T(imm, 11)
	c := bi.Ts(imm, 5, 6)
	d := bi.Ts(imm, 1, 4)
	return a<<31 | c<<25 | rs2<<20 | rs1<<15 | f3<<12 | d<<8 | b<<7 | op<<2 | 3
}
