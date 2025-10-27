package instr

import "github.com/temnok/rv/state"

func Bgeu(cpu *state.State, rs1, rs2, imm int) {
	a := cpu.Xuint(cpu.X[rs1])
	b := cpu.Xuint(cpu.X[rs2])

	c := a >= b

	cpu.PCAddIf(c, imm)
}
