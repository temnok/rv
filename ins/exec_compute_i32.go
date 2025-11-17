package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execComputeI32(cpu *state.CPU, op Op) {
	imm := imm.I(op.code())

	if cpu.LenIs64() {
		switch op.f3() {
		case 0b_000:
			addiw(cpu, op)
		case 0b_001:
			if imm < 32 {
				slliw(cpu, op)
			}
		case 0b_101:
			if imm < 32 {
				srliw(cpu, op)
			} else if imm&^0b0100000_00000 < 32 {
				sraiw(cpu, op)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		illegal(cpu, op)
	}
}

func computeI32(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.rs1()])
	b := int32(imm.I(op.code()))

	c := int(f(a, b))

	cpu.Xset(op.rd(), c)
}
