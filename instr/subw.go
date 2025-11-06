package instr

import "github.com/temnok/rv/state"

func Subw(cpu *state.CPU, op Op) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(cpu.X[op.Rs2()])

	c := int(a - b)

	cpu.Xset(op.Rd(), c)
}
