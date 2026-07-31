package inst

import "github.com/temnok/rv/state"

func (ctx *context) SLT(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	if ctx.X[rs1] < ctx.X[rs2] {
		ctx.Update.Xval = 1
	} else {
		ctx.Update.Xval = 0
	}
}
