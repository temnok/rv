package instr

import "github.com/temnok/rv/state"

func Slt(cpu *state.State, rd, rs1, rs2 int) {
	if cpu.X[rs1] < cpu.X[rs2] {
		cpu.Xset(rd, 1)
	} else {
		cpu.Xset(rd, 0)
	}
}
