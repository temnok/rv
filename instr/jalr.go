package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Jalr(cpu *state.State, op Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	savedPC := cpu.Update.PC
	newPC := (cpu.X[rs1] + imm) &^ 1

	cpu.Xset(rd, savedPC)
	cpu.Update.PC = cpu.Xint(newPC)
}
