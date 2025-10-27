package instr

import "github.com/temnok/rv/state"

func Sllw(cpu *state.State, rd, rs1, rs2 int) {
	a := int32(cpu.X[rs1])
	b := int32(cpu.X[rs2]) & 31

	c := int(a << b)

	cpu.Xset(rd, c)
}
