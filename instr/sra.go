package instr

import "github.com/temnok/rv/state"

func Sra(cpu *state.State, rd, rs1, rs2 int) {
	cpu.XRegSet(rd, cpu.X[rs1]>>(cpu.X[rs2]&(cpu.XLen-1)))
}
