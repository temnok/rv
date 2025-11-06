package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Srai(cpu *state.CPU, op Op) {
	a := cpu.X[op.Rs1()]
	b := imm.I(op.Code()) & cpu.Xmask()

	c := a >> b

	cpu.Xset(op.Rd(), c)
}
