package inst

import "github.com/temnok/rv/state"

func (ctx *context) SUBW(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a := int32(ctx.X[rs1])
	b := int32(ctx.X[rs2])

	ctx.Update.Xval = int(a - b)
}
