package imm

import "github.com/temnok/rv/bit"

func CL(op int) int {
	a := bit.Get(op, 5)
	b := bit.GetN(op, 10, 3)
	c := bit.Get(op, 6)

	return a<<6 | b<<3 | c<<2
}
