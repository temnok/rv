package instr

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Ebreak(cpu *state.CPU, op Op) {
	trap.EnterWithoutTval(cpu, trap.Breakpoint)
}
