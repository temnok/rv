package imm

import "github.com/temnok/rv/bit"

func CSS3(op int) int {
	a := bit.GetN(op, 7, 3)
	b := bit.GetN(op, 10, 3)

	return a<<6 | b<<3
}
