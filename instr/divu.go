package instr

import (
	"github.com/temnok/rv/state"
)

func Divu(cpu *state.CPU, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(cpu.X[op.Rs2()])

	c := -1
	if b != 0 {
		c = int(a / b)
	}

	cpu.Xset(op.Rd(), c)
}
