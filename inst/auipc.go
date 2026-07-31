package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) AUIPC(rd, imm int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	ctx.Update.Xval = ctx.PC + imm
}
