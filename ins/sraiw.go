package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Sraiw(cpu *state.CPU, op Op) {
	a := int32(cpu.X[op.Rs1()])
	b := int32(imm.I(op.Code())) & 31

	c := int(a >> b)

	cpu.Xset(op.Rd(), c)
}
