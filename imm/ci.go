package imm

import "github.com/temnok/rv/bit"

func CI(op int) int {
	a := bit.Get(op, 12)
	b := bit.GetN(op, 2, 5)

	return -a<<5 | b
}
