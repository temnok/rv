package instr

import "github.com/temnok/rv/state"

func Bge(cpu *state.State, rs1, rs2, imm int) {
	a := cpu.X[rs1]
	b := cpu.X[rs2]

	c := a >= b

	if c {
		cpu.PCAdd(imm)
	}
}
