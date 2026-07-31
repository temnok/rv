package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) SRAIW(rd, rs1, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	ctx.Update.Xval = int(int32(ctx.X[rs1]) >> int32(imm&0x1F))
}
