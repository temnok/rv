package ins

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ebreak(cpu *state.CPU, op Op) {
	trap.EnterWithoutTval(cpu, trap.Breakpoint)
}
