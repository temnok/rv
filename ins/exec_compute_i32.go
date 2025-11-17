package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execComputeI32(cpu *state.CPU, op Op) {
	imm := imm.I(op.Code())

	if cpu.LenIs64() {
		switch op.F3() {
		case 0b_000:
			Addiw(cpu, op)
		case 0b_001:
			if imm < 32 {
				Slliw(cpu, op)
			}
		case 0b_101:
			if imm < 32 {
				Srliw(cpu, op)
			} else if imm&^0b0100000_00000 < 32 {
				Sraiw(cpu, op)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		illegal(cpu, op)
	}
}

func computeI32(cpu *state.CPU, op Op, f func(a, b int32) int32) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(imm.I(op.Code()))

	c := int(f(a, b))

	cpu.Xset(op.Rd(), c)
}
