package imm

import "github.com/temnok/rv/bit"

func CI3(op int) int {
	a := bit.GetN(op, 2, 3)
	b := bit.Get(op, 12)
	c := bit.GetN(op, 5, 2)

	return a<<6 | b<<5 | c<<3
}
