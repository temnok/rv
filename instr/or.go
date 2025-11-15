package instr

import "github.com/temnok/rv/state"

func or(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		return a | b
	})
}
