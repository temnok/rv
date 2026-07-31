package inst

import "github.com/temnok/rv/state"

func (ctx *context) DIVUW(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a, b := uint32(ctx.X[rs1]), uint32(ctx.X[rs2])

	if b != 0 {
		ctx.Update.Xval = int(int32(a / b))
	} else {
		ctx.Update.Xval = -1
	}
}
