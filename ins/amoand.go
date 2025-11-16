package ins

import (
	"github.com/temnok/rv/state"
)

func amoand(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		*val &= *old

		return true
	})
}
