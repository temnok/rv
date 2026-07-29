package inst

import (
	"github.com/temnok/rv/state"
)

func divu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if b == 0 {
			return -1
		}

		return int(uint(a) / uint(b))
	})
}
