package instr

import "github.com/temnok/rv/state"

func sll(cpu *state.CPU, op Op) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()] & cpu.Xmask()

	c := a << b

	cpu.Xset(op.Rd(), c)
}
