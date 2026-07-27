package imm

import "github.com/temnok/rv/bit"

func J(op int) int {
	a := bit.Get(op, 31)
	b := bit.GetN(op, 12, 8)
	c := bit.Get(op, 20)
	d := bit.GetN(op, 21, 10)

	return -a<<20 | b<<12 | c<<11 | d<<1
}
