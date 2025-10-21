package instr

import (
	"github.com/temnok/rv/state"
)

func Divu(cpu *state.State, rd, rs1, rs2 int) {
	a, b, c := cpu.X[rs1], cpu.X[rs2], 0

	if b != 0 {
		c = int(cpu.Xuint(a) / cpu.Xuint(b))
	} else {
		c = -1
	}

	cpu.Xset(rd, c)
}
