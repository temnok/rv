package ins

import (
	"github.com/temnok/rv/state"
)

func amoswap(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		return true
	})
}
