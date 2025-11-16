package ins

import (
	"github.com/temnok/rv/state"
)

func amomin(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		if *old < *val {
			*val = *old
		}

		return true
	})
}
