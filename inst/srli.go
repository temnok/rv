package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SRLI(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a := uint(ctx.X[rs1])
	b := uint(imm) & 0x3F

	ctx.Update.Xval = int(a >> b)
}
