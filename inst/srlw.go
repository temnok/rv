package inst

import "github.com/temnok/rv/state"

func srlw(cpu *state.CPU, op Op) {
	computeR32(cpu, op, func(a, b int32) int32 {
		b &= 31

		return int32(uint32(a) >> uint32(b))
	})
}
