package encode

import "github.com/temnok/rv/bit"

func J(imm, rd, op int) int {
	imm10_1 := bit.GetN(imm, 1, 10)
	imm11 := bit.Get(imm, 11)
	imm19_12 := bit.GetN(imm, 12, 8)
	imm20 := bit.Get(imm, 20)

	return imm20<<31 | imm10_1<<21 | imm11<<20 | imm19_12<<12 | rd<<7 | op<<2 | 3
}
