package instr

import "github.com/temnok/rv/state"

func Remuw(cpu *state.State, op Op) {
	a := uint32(cpu.X[op.Rs1()])
	b := uint32(cpu.X[op.Rs2()])

	c := int(int32(a))
	if b != 0 {
		c = int(int32(a % b))
	}

	cpu.Xset(op.Rd(), c)
}
