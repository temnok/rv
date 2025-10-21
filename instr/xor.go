package instr

import "github.com/temnok/rv/state"

func Xor(cpu *state.State, rd, rs1, rs2 int) {
	cpu.Xset(rd, cpu.X[rs1]^cpu.X[rs2])
}
