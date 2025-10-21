package instr

import "github.com/temnok/rv/state"

func Slti(cpu *state.State, rd, rs1, imm int) {
	if cpu.X[rs1] < imm {
		cpu.Xset(rd, 1)
	} else {
		cpu.Xset(rd, 0)
	}
}
