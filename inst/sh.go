package inst

import (
	"github.com/temnok/rv/state"
)

func sh(cpu *state.CPU, op Op) {
	store(cpu, op, 2)
}
