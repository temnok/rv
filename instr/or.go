package instr

import "github.com/temnok/rv/state"

func Or(cpu *state.State, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := a | b

	cpu.Xset(op.Rd(), c)
}
