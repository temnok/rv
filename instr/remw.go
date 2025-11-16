package instr

import "github.com/temnok/rv/state"

func remw(cpu *state.CPU, op Op) {
	computeR64(cpu, op, func(a, b int32) int32 {
		if b == 0 {
			return a
		}

		return a % b
	})
}
