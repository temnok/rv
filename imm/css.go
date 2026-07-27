package imm

import "github.com/temnok/rv/bit"

func CSS(op int) int {
	a := bit.GetN(op, 7, 2)
	b := bit.GetN(op, 9, 4)

	return a<<6 | b<<2
}
