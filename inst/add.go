package inst

import "github.com/temnok/rv/state"

func (ctx *context) ADD(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	ctx.Update.Xval = ctx.X[rs1] + ctx.X[rs2]
}
