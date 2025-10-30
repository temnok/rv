package instr

import "github.com/temnok/rv/state"

func Slt(cpu *state.State, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := a < b

	cpu.XsetBool(op.Rd(), c)
}
