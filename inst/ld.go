package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) LD(rd, rs1, offset int) {
	if val, ok := memRead(ctx, ctx.X[rs1]+offset, 8); ok {
		ctx.Update.Targets = state.UpdateXreg

		ctx.Update.Xreg = rd
		ctx.Update.Xval = val
	}
}
