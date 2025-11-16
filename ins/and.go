package ins

import "github.com/temnok/rv/state"

func and(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		return a & b
	})
}
