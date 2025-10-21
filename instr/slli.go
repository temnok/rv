package instr

import "github.com/temnok/rv/state"

func Slli(cpu *state.State, rd, rs1, imm int) {
	if imm < cpu.XLen {
		cpu.Xset(rd, cpu.X[rs1]<<imm)
	}
}
