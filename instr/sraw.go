package instr

import "github.com/temnok/rv/state"

func Sraw(cpu *state.State, rd, rs1, rs2 int) {
	shamt := int32(cpu.X[rs2]) & 31
	cpu.Xset(rd, int(int32(cpu.X[rs1])>>shamt))
}
