package instr

import "github.com/temnok/rv/state"

func Divw(cpu *state.State, op Op) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(cpu.X[op.Rs2()])

	c := -1
	if b != 0 {
		c = int(a / b)
	}

	cpu.Xset(op.Rd(), c)
}
