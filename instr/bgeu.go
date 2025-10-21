package instr

import "github.com/temnok/rv/state"

func Bgeu(cpu *state.State, rs1, rs2, imm int) {
	if cpu.Xuint(cpu.X[rs1]) >= cpu.Xuint(cpu.X[rs2]) {
		cpu.PCAdd(imm)
	}
}
