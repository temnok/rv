package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Addi(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		return a + b
	})
}

func computeI(cpu *state.CPU, op Op, f func(a, b int) int) {
	a := cpu.X[op.Rs1()]
	b := imm.I(op.Code())

	c := f(a, b)

	cpu.Xset(op.Rd(), c)
}
