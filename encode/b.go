package encode

import "github.com/temnok/rv/bit"

func B(imm, rs2, rs1, f3, op int) int {
	imm4_1 := bit.GetN(imm, 1, 4)
	imm10_5 := bit.GetN(imm, 5, 6)
	imm11 := bit.Get(imm, 11)
	imm12 := bit.Get(imm, 12)

	return imm12<<31 | imm10_5<<25 | rs2<<20 | rs1<<15 | f3<<12 | imm4_1<<8 | imm11<<7 | op<<2 | 3
}
