package instr

import "github.com/temnok/rv/state"

func Sltu(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(cpu.X[rs2])

	c := a < b

	cpu.XsetBool(rd, c)
}
