package imm

import "github.com/temnok/rv/bit"

func CI2(op int) int {
	a := bit.GetN(op, 2, 2)
	b := bit.Get(op, 12)
	c := bit.GetN(op, 4, 3)

	return a<<6 | b<<5 | c<<2
}
