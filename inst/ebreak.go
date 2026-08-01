package inst

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ebreak(cpu *state.CPU, op Op) {
	(*context)(cpu).EBREAK()
}

func (ctx *context) EBREAK() {
	trap.Enter((*state.CPU)(ctx), trap.Breakpoint, 0)
}
