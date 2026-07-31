package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) DIVU(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	if ctx.X[rs2] != 0 {
		ctx.Update.Xval = int(uint(ctx.X[rs1]) / uint(ctx.X[rs2]))
	} else {
		ctx.Update.Xval = -1
	}
}
