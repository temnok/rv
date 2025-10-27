package instr

import "github.com/temnok/rv/state"

func Srai(cpu *state.State, rd, rs1, imm int) {
	a := cpu.X[rs1]
	b := imm

	c := a >> b

	cpu.Xset(rd, c)
}
