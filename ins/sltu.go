package ins

import "github.com/temnok/rv/state"

func sltu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if uint(a) < uint(b) {
			return 1
		}

		return 0
	})
}
