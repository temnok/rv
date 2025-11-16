package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func csrRS(cpu *state.CPU, op Op) (int, int) {
	r := bi.Ts(imm.I(op.Code()), 0, 12)

	s := op.Rs1()
	if (op.F3() & 4) == 0 {
		s = cpu.X[s]
	}

	return r, s
}
