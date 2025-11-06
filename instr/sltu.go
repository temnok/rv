package instr

import "github.com/temnok/rv/state"

func Sltu(cpu *state.CPU, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(cpu.X[op.Rs2()])

	c := a < b

	cpu.XsetBool(op.Rd(), c)
}
