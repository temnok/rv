package instr

import "github.com/temnok/rv/state"

func Slliw(cpu *state.State, rd, rs1, imm int) {
	// TODO
	if imm < 32 {
		a := int32(cpu.X[rs1])
		b := int32(imm)

		c := int(a << b)

		cpu.Xset(rd, c)
	}
}
