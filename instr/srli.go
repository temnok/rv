package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Srli(cpu *state.CPU, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(imm.I(op.Code()))

	c := int(a >> b)

	cpu.Xset(op.Rd(), c)
}
