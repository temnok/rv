package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func jalr(cpu *state.CPU, op Op) {
	imm := imm.I(op.code())

	savedPC := cpu.Update.PC
	newPC := (cpu.X[op.rs1()] + imm) &^ 1

	cpu.Xset(op.rd(), savedPC)
	cpu.Update.PC = cpu.Int(newPC)
}
