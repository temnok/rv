package ins

import (
	"github.com/temnok/rv/state"
)

func amominu(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		if cpu.Uint(*old) < cpu.Uint(*val) {
			*val = *old
		}

		return true
	})
}
