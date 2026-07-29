package inst

import (
	"github.com/temnok/rv/state"
)

func sraiw(cpu *state.CPU, op Op) {
	computeI32(cpu, op, func(a, b int32) int32 {
		b &= 31

		return a >> b
	})
}
