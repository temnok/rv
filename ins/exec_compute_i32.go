package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execComputeI32(cpu *state.CPU, op Op) {
	imm := imm.I(op.code())

	ins := illegal

	switch op.f3() {
	case 0:
		ins = addiw
	case 1:
		if imm < 32 {
			ins = slliw
		}
	case 5:
		if imm < 32 {
			ins = srliw
		} else if imm&^0b0100000_00000 < 32 {
			ins = sraiw
		}
	}

	ins(cpu, op)
}

func computeI32(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.rs1()])
	b := int32(imm.I(op.code()))

	c := int(f(a, b))

	cpu.Xset(op.rd(), c)
}
