package ins

import (
	"github.com/temnok/rv/state"
)

func sltiu(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		if uint(a) < uint(b) {
			return 1
		}

		return 0
	})
}
