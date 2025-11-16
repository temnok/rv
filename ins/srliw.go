package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Srliw(cpu *state.CPU, op Op) {
	a := uint32(cpu.X[op.Rs1()])
	b := uint32(imm.I(op.Code())) & 31

	c := int(int32(a >> b))

	cpu.Xset(op.Rd(), c)
}
