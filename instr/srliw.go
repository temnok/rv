package instr

import "github.com/temnok/rv/state"

func Srliw(cpu *state.State, rd, rs1, imm int) {
	a := uint32(cpu.X[rs1])
	b := uint32(imm)

	c := int(int32(a >> b))

	cpu.Xset(rd, c)
}
