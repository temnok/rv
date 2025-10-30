package instr

import "github.com/temnok/rv/state"

func Mulw(cpu *state.State, op Op) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(cpu.X[op.Rs2()])

	c := int(a * b)

	cpu.Xset(op.Rd(), c)
}
