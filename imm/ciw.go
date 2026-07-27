package imm

import (
	"github.com/temnok/rv/bit"
)

func CIW(op int) int {
	a := bit.GetN(op, 7, 4)
	b := bit.GetN(op, 11, 2)
	c := bit.Get(op, 5)
	d := bit.Get(op, 6)

	return a<<6 | b<<4 | c<<3 | d<<2
}
