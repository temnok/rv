package instr

import "github.com/temnok/rv/state"

func Sraiw(cpu *state.State, rd, rs1, imm int) {
	a := int32(cpu.X[rs1])
	b := int32(imm)

	c := int(a >> b)

	cpu.Xset(rd, c)
}
