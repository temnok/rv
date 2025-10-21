package instr

import "github.com/temnok/rv/state"

func Srl(cpu *state.State, rd, rs1, rs2 int) {
	cpu.Xset(rd, int(cpu.Xuint(cpu.X[rs1])>>cpu.Xuint(cpu.X[rs2]&(cpu.XLen-1))))
}
