package instr

import "github.com/temnok/rv/state"

func Srli(cpu *state.State, rd, rs1, imm int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(imm)

	c := int(a >> b)

	cpu.Xset(rd, c)
}
