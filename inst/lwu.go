package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) LWU(rd, rs1, offset int) {
	if val, ok := memRead(ctx, ctx.X[rs1]+offset, 4); ok {
		ctx.Update.Targets |= state.UpdateXreg // device read can set other targets

		ctx.Update.Xreg = rd
		ctx.Update.Xval = int(uint32(val))
	}
}
