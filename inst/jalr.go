package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func jalr(cpu *state.CPU, op Op) {
	cpu.Update.Targets = state.UpdateXreg | state.UpdatePC

	cpu.Update.Xreg = op.rd()
	cpu.Update.Xval = cpu.Update.FollowPC

	cpu.Update.PC = (cpu.X[op.rs1()] + imm.I(op.code())) &^ 1
}
