package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SLTI(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	if ctx.X[rs1] < imm {
		ctx.Update.Xval = 1
	} else {
		ctx.Update.Xval = 0
	}
}
