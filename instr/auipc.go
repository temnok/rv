package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Auipc(cpu *state.State, op Op) {
	imm, rd := imm.U(op.Code()), op.Rd()

	newPC := cpu.PC + imm

	cpu.Xset(rd, newPC)
}
