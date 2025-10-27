package instr

import "github.com/temnok/rv/state"

func Sll(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.X[rs1]
	b := cpu.X[rs2] & (cpu.XLen - 1)

	c := a << b

	cpu.Xset(rd, c)
}
