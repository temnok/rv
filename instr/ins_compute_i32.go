package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func computeI32(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(imm.I(op.Code()))

	c := int(f(a, b))

	cpu.Xset(op.Rd(), c)
}
