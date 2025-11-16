package ins

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Mret(cpu *state.CPU, op Op) {
	trap.Exit(cpu, state.PrivM)
}
