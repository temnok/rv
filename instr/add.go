package instr

import "github.com/temnok/rv/state"

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#norm:add_op
func Add(cpu *state.State, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := a + b

	cpu.Xset(op.Rd(), c)
}
