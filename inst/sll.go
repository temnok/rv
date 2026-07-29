package inst

import (
	"github.com/temnok/rv/state"
)

func sll(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= 63

		return a << b
	})
}
