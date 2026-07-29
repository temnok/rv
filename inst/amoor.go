package inst

import (
	"github.com/temnok/rv/state"
)

func amoor(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		*val |= *old

		return true
	})
}
