package instr

import "github.com/temnok/rv/state"

func Srlw(cpu *state.State, rd, rs1, rs2 int) {
	a := uint32(cpu.X[rs1])
	b := uint32(cpu.X[rs2]) & 31

	c := int(int32(a >> b))

	cpu.Xset(rd, c)
}
