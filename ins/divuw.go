package ins

import "github.com/temnok/rv/state"

func divuw(cpu *state.CPU, op Op) {
	computeR32(cpu, op, func(a, b int32) int32 {
		if b == 0 {
			return -1
		}

		return int32(uint32(a) / uint32(b))
	})
}
