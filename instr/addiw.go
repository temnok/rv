package instr

import "github.com/temnok/rv/state"

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#norm:addiw_op
func Addiw(cpu *state.State, rd, rs1, imm int) {
	a := int32(cpu.X[rs1])
	b := int32(imm)

	c := int(a + b)

	cpu.Xset(rd, c)
}
