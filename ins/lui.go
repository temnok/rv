package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func lui(cpu *state.CPU, op Op) {
	imm := imm.U(op.code())

	cpu.Xset(op.rd(), imm)
}
