package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Ori(cpu *state.State, op Op) {
	a := cpu.X[op.Rs1()]
	b := imm.I(op.Code())

	c := a | b

	cpu.Xset(op.Rd(), c)
}
