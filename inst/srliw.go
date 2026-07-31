package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SRLIW(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a := uint32(ctx.X[rs1])
	b := uint32(imm) & 0x1F

	ctx.Update.Xval = int(int32(a >> b))
}
