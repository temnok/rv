package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) JALR(rd, rs1, offset int) {
	ctx.Update.Targets = state.UpdateXreg | state.UpdatePC

	ctx.Update.Xreg = rd
	ctx.Update.Xval = ctx.Update.FollowPC

	ctx.Update.PC = (ctx.X[rs1] + offset) &^ 1
}
