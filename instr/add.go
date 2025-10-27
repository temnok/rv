package instr

import "github.com/temnok/rv/state"

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#norm:add_op
func Add(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.X[rs1]
	b := cpu.X[rs2]

	c := a + b

	cpu.Xset(rd, c)
}
