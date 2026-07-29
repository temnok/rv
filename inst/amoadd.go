package inst

import (
	"github.com/temnok/rv/state"
)

func amoadd(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		*val += *old

		return true
	})
}
