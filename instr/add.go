package instr

import "github.com/temnok/rv/state"

func Add(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		return a + b
	})
}

func computeR(cpu *state.CPU, op Op, f func(a, b int) int) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := f(a, b)

	cpu.Xset(op.Rd(), c)
}
