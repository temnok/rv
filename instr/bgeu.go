package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Bgeu(cpu *state.State, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(cpu.X[op.Rs2()])

	c := a >= b

	cpu.PCAddIf(c, imm.B(op.Code()))
}
