package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Sltiu(cpu *state.State, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(imm.I(op.Code()))

	c := a < b

	cpu.XsetBool(op.Rd(), c)
}
