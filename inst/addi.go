package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) ADDI(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	ctx.Update.Xval = ctx.X[rs1] + imm
}
