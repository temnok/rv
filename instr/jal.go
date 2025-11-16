package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Jal(cpu *state.CPU, op Op) {
	imm := imm.J(op.Code())

	savedPC := cpu.Update.PC
	newPC := cpu.Int(cpu.PC + imm)

	cpu.Xset(op.Rd(), savedPC)
	cpu.Update.PC = newPC
}
