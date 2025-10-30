package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Jalr(cpu *state.State, op Op) {
	imm := imm.I(op.Code())

	savedPC := cpu.Update.PC
	newPC := (cpu.X[op.Rs1()] + imm) &^ 1

	cpu.Xset(op.Rd(), savedPC)
	cpu.Update.PC = cpu.Xint(newPC)
}
