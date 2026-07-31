package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) JAL(rd, offset int) {
	ctx.Update.Targets = state.UpdateXreg | state.UpdatePC

	ctx.Update.Xreg = rd
	ctx.Update.Xval = ctx.Update.FollowPC

	ctx.Update.PC = ctx.PC + offset
}
