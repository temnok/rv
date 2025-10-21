package instr

import "github.com/temnok/rv/state"

func Srlw(cpu *state.State, rd, rs1, rs2 int) {
	shamt := uint32(cpu.X[rs2]) & 31
	cpu.XRegSet(rd, int(int32(uint32(cpu.X[rs1])>>shamt)))
}
