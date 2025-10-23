package instr

import (
	"github.com/temnok/rv/state"
)

func Divu(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(cpu.X[rs2])

	c := -1
	if b != 0 {
		c = int(a / b)
	}

	cpu.Xset(rd, c)
}
