package imm

import "github.com/temnok/rv/bit"

func U(op int) int {
	a := bit.Get(op, 31)
	b := bit.GetN(op, 12, 19)

	return -a<<31 | b<<12
}
