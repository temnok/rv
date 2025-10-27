package instr

import "github.com/temnok/rv/state"

func Sltiu(cpu *state.State, rd, rs1, imm int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(imm)

	c := a < b

	cpu.XsetBool(rd, c)
}
