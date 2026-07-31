package inst

import (
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

func (ctx *context) SH(rs2, rs1, offset int) {
	mem.Write((*state.CPU)(ctx), ctx.X[rs1]+offset, 2, ctx.X[rs2])
}
