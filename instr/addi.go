package instr

import "github.com/temnok/rv/state"

func Addi(cpu *state.State, rd, rs1, imm int) {
	cpu.XRegSet(rd, cpu.X[rs1]+imm)
}
