package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func jal(cpu *state.CPU, op Op) {
	imm := imm.J(op.code())

	savedPC := cpu.Update.PC
	newPC := cpu.PC + imm

	cpu.Xset(op.rd(), savedPC)
	cpu.Update.PC = newPC
}
