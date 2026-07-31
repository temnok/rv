package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SLTIU(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	if uint(ctx.X[rs1]) < uint(imm) {
		ctx.Update.Xval = 1
	} else {
		ctx.Update.Xval = 0
	}
}
