package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Lui(cpu *state.State, op Op) {
	imm, rd := imm.U(op.Code()), op.Rd()

	cpu.Xset(rd, imm)
}
