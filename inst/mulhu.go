package inst

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func (ctx *context) MULHU(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a, b := ctx.X[rs1], ctx.X[rs2]
	hi, _ := bits.Mul64(uint64(a), uint64(b))

	ctx.Update.Xval = int(hi)
}
