package instr

import (
	"github.com/temnok/rv/state"
)

func Rem(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.X[rs1]
	b := cpu.X[rs2]

	c := a
	if b != 0 {
		c = a % b
	}

	cpu.Xset(rd, c)
}
