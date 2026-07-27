package imm

import "github.com/temnok/rv/bit"

func CJ(op int) int {
	a := bit.Get(op, 12)
	b := bit.Get(op, 8)
	c := bit.GetN(op, 9, 2)
	d := bit.Get(op, 6)
	e := bit.Get(op, 7)
	f := bit.Get(op, 2)
	g := bit.Get(op, 11)
	h := bit.GetN(op, 3, 3)

	return -a<<11 | b<<10 | c<<8 | d<<7 | e<<6 | f<<5 | g<<4 | h<<1
}
