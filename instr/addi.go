package instr

import "github.com/temnok/rv/state"

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#norm:addi_op
func Addi(cpu *state.State, rd, rs1, imm int) {
	a := cpu.X[rs1]
	b := imm

	c := a + b

	cpu.Xset(rd, c)
}
