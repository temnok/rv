package inst

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ecall(cpu *state.CPU, op Op) {
	(*context)(cpu).ECALL()
}

func (ctx *context) ECALL() {
	trap.Enter((*state.CPU)(ctx), trap.EnvironmentCallFromUMode+ctx.CSR.Priv, 0)
}
