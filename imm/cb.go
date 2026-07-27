package imm

import "github.com/temnok/rv/bit"

func CB(op int) int {
	a := bit.Get(op, 12)
	b := bit.GetN(op, 5, 2)
	c := bit.Get(op, 2)
	d := bit.GetN(op, 10, 2)
	e := bit.GetN(op, 3, 2)

	return -a<<8 | b<<6 | c<<5 | d<<3 | e<<1
}
