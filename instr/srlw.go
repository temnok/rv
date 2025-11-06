package instr

import "github.com/temnok/rv/state"

func Srlw(cpu *state.CPU, op Op) {
	a := uint32(cpu.X[op.Rs1()])
	b := uint32(cpu.X[op.Rs2()]) & 31

	c := int(int32(a >> b))

	cpu.Xset(op.Rd(), c)
}
