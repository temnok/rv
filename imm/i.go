package imm

import "github.com/temnok/rv/bit"

func I(op int) int {
	a := bit.Get(op, 31)
	b := bit.GetN(op, 20, 11)

	return -a<<11 | b
}
