package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) REM(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	if ctx.X[rs2] != 0 {
		ctx.Update.Xval = ctx.X[rs1] % ctx.X[rs2]
	} else {
		ctx.Update.Xval = ctx.X[rs1]
	}
}
