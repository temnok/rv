package instr

import "github.com/temnok/rv/state"

func Srl(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(cpu.X[rs2] & cpu.Xmask())

	c := int(a >> b)

	cpu.Xset(rd, c)
}
