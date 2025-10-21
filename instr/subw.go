package instr

import "github.com/temnok/rv/state"

func Subw(cpu *state.State, rd, rs1, rs2 int) {
	cpu.Xset(rd, int(int32(cpu.X[rs1])-int32(cpu.X[rs2])))
}
