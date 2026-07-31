package inst

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func (ctx *context) MULH(rd, rs1, rs2 int) {
	ctx.Update.Targets = state.UpdateXreg

	ctx.Update.Xreg = rd

	a, b := ctx.X[rs1], ctx.X[rs2]
	hi, _ := bits.Mul64(uint64(a), uint64(b))
	s1 := (a >> 63) & b
	s2 := (b >> 63) & a

	ctx.Update.Xval = int(hi) - s1 - s2
}
