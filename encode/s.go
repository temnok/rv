package encode

import "github.com/temnok/rv/bi"

func S(imm, rs2, rs1, f3, op int) int {
	a := bi.Ts(imm, 5, 7)
	b := bi.Ts(imm, 0, 5)
	return a<<25 | rs2<<20 | rs1<<15 | f3<<12 | b<<7 | op<<2 | 3
}
