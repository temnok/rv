package inst

import (
	"github.com/temnok/rv/state"
)

func (ctx *context) BGEU(rs1, rs2, offset int) {
	if uint(ctx.X[rs1]) >= uint(ctx.X[rs2]) {
		ctx.Update.Targets = state.UpdatePC
		ctx.Update.PC = ctx.PC + offset
	}
}
