package instr

import "github.com/temnok/rv/state"

func Sltiu(cpu *state.State, rd, rs1, imm int) {
	if cpu.Xuint(cpu.X[rs1]) < cpu.Xuint(imm) {
		cpu.XRegSet(rd, 1)
	} else {
		cpu.XRegSet(rd, 0)
	}
}
