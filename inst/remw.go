package inst

import "github.com/temnok/rv/state"

func (ctx *context) REMW(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a, b := int32(ctx.X[rs1]), int32(ctx.X[rs2])

	if b != 0 {
		ctx.Update.Xval = int(a % b)
	} else {
		ctx.Update.Xval = int(a)
	}
}
