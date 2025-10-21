package instr

import "github.com/temnok/rv/state"

func Sltiu(cpu *state.State, rd, rs1, imm int) {
	if cpu.Xuint(cpu.X[rs1]) < cpu.Xuint(imm) {
		cpu.Xset(rd, 1)
	} else {
		cpu.Xset(rd, 0)
	}
}
