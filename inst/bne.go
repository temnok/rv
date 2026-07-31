package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) BNE(rs1, rs2, offset int) {
	if ctx.X[rs1] != ctx.X[rs2] {
		ctx.Update.Targets = state.UpdatePC
		ctx.Update.PC = ctx.PC + offset
	}
}
