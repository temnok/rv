package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Slliw(cpu *state.CPU, op Op) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(imm.I(op.Code()))

	c := int(a << b)

	cpu.Xset(op.Rd(), c)
}
