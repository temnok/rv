package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SRL(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a := uint(ctx.X[rs1])
	b := uint(ctx.X[rs2]) & 0x3F

	ctx.Update.Xval = int(a >> b)
}
