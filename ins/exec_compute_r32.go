package ins

import (
	"github.com/temnok/rv/state"
)

func execComputeR32(cpu *state.CPU, op Op) {
	f3, f7 := op.f3(), op.f7()

	ins := illegal

	switch f7 {
	case 0:
		switch f3 {
		case 0:
			ins = addw
		case 1:
			ins = sllw
		case 5:
			ins = srlw
		}
	case 1:
		ins = computeM32
	case 1 << 5:
		switch f3 {
		case 0:
			ins = subw
		case 5:
			ins = sraw
		}
	}

	ins(cpu, op)
}

func computeR32(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.rs1()])
	b := int32(cpu.X[op.rs2()])

	c := int(f(a, b))

	cpu.Xset(op.rd(), c)
}
