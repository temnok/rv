package instr

import "github.com/temnok/rv/state"

func And(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		return a & b
	})
}
