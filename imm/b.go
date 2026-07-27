package imm

import "github.com/temnok/rv/bit"

func B(op int) int {
	a := bit.Get(op, 31)
	b := bit.Get(op, 7)
	c := bit.GetN(op, 25, 6)
	d := bit.GetN(op, 8, 4)

	return -a<<12 | b<<11 | c<<5 | d<<1
}
