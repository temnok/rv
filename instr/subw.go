package instr

import "github.com/temnok/rv/state"

func subw(cpu *state.CPU, op Op) {
	computeR32(cpu, op, func(a, b int32) int32 {
		return a - b
	})
}
