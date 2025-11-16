package ins

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Sret(cpu *state.CPU, op Op) {
	trap.Exit(cpu, state.PrivS)
}
