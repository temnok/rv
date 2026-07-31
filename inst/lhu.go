package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) LHU(rd, rs1, offset int) {
	if val, ok := memRead(ctx, ctx.X[rs1]+offset, 2); ok {
		ctx.Update.Targets = state.UpdateXreg

		ctx.Update.Xreg = rd
		ctx.Update.Xval = int(uint16(val))
	}
}
