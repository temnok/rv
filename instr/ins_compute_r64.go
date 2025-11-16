package instr

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func ComputeR64(cpu *state.CPU, op Op) {
	if !cpu.LenIs64() {
		illegal(cpu, op)
		return
	}

	f7, f3 := op.F7(), op.F3()

	if f7 == 1 {
		computeM64(cpu, op)
		return
	}

	if f7&^0b0100000 != 0 {
		illegal(cpu, op)
		return
	}

	switch bi.T(f7, 5)<<3 | f3 {
	case 0b_000:
		Addw(cpu, op)
	case 0b_1_000:
		Subw(cpu, op)
	case 0b_001:
		Sllw(cpu, op)
	case 0b_101:
		Srlw(cpu, op)
	case 0b_1_101:
		Sraw(cpu, op)
	default:
		illegal(cpu, op)
	}
}

func computeR64(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(cpu.X[op.Rs2()])

	c := int(f(a, b))

	cpu.Xset(op.Rd(), c)
}
