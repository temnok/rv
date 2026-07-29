package inst

import "github.com/temnok/rv/state"

func sllw(cpu *state.CPU, op Op) {
	computeR32(cpu, op, func(a, b int32) int32 {
		b &= 31

		return a << b
	})
}
