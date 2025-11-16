package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Jalr(cpu *state.CPU, op Op) {
	imm := imm.I(op.Code())

	savedPC := cpu.Update.PC
	newPC := (cpu.X[op.Rs1()] + imm) &^ 1

	cpu.Xset(op.Rd(), savedPC)
	cpu.Update.PC = cpu.Int(newPC)
}
