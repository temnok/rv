package instr

import "github.com/temnok/rv/state"

func Bge(cpu *state.State, rs1, rs2, imm int) {
	if cpu.X[rs1] >= cpu.X[rs2] {
		cpu.PCAdd(imm)
	}
}
