package instr

import "github.com/temnok/rv/state"

func slt(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if a < b {
			return 1
		}

		return 0
	})
}
