package imm

import "github.com/temnok/rv/bit"

func S(op int) int {
	a := bit.Get(op, 31)
	b := bit.GetN(op, 25, 6)
	c := bit.GetN(op, 7, 5)

	return -a<<11 | b<<5 | c
}
