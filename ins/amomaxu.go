package ins

import (
	"github.com/temnok/rv/state"
)

func amomaxu(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		if uint(*old) > uint(*val) {
			*val = *old
		}

		return true
	})
}
