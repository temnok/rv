package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Jal(cpu *state.State, op Op) {
	imm, rd := imm.J(op.Code()), op.Rd()

	savedPC := cpu.Update.PC
	newPC := cpu.Xint(cpu.PC + imm)

	cpu.Xset(rd, savedPC)
	cpu.Update.PC = newPC
}
