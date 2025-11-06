package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Bge(cpu *state.CPU, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := a >= b

	cpu.PCAddIf(c, imm.B(op.Code()))
}
