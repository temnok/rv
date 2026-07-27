package encode

import "github.com/temnok/rv/bit"

func S(imm, rs2, rs1, f3, op int) int {
	imm4_0 := bit.GetN(imm, 0, 5)
	imm11_5 := bit.GetN(imm, 5, 7)

	return imm11_5<<25 | rs2<<20 | rs1<<15 | f3<<12 | imm4_0<<7 | op<<2 | 3
}
