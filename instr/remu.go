package instr

import (
	"github.com/temnok/rv/state"
)

func Remu(cpu *state.CPU, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := a
	if b != 0 {
		c = int(cpu.Xuint(a) % cpu.Xuint(b))
	}

	cpu.Xset(op.Rd(), c)
}
