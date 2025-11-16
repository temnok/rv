package instr

import (
	"github.com/temnok/rv/state"
)

func slti(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		if a < b {
			return 1
		}

		return 0
	})
}
