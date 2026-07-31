package inst

import "github.com/temnok/rv/state"

func (ctx *context) SRLW(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd
	ctx.Update.Xval = int(int32(uint32(ctx.X[rs1]) >> uint32(ctx.X[rs2]&0x1F)))
}
