package instr

import "github.com/temnok/rv/state"

func sltu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if cpu.Xuint(a) < cpu.Xuint(b) {
			return 1
		}

		return 0
	})
}
