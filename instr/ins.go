package instr

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

type instr = func(*state.CPU, Op)

func illegal(cpu *state.CPU, op Op) {
	trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
}
