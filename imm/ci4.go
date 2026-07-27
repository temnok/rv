package imm

import "github.com/temnok/rv/bit"

func CI4(op int) int {
	a := bit.Get(op, 12)
	b := bit.GetN(op, 3, 2)
	c := bit.Get(op, 5)
	d := bit.Get(op, 2)
	e := bit.Get(op, 6)

	return -a<<9 | b<<7 | c<<6 | d<<5 | e<<4
}
