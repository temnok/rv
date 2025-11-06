package instr

import "github.com/temnok/rv/state"

func Srl(cpu *state.CPU, op Op) {
	a := cpu.Xuint(cpu.X[op.Rs1()])
	b := cpu.Xuint(cpu.X[op.Rs2()] & cpu.Xmask())

	c := int(a >> b)

	cpu.Xset(op.Rd(), c)
}
