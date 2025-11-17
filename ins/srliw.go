package ins

import (
	"github.com/temnok/rv/state"
)

func srliw(cpu *state.CPU, op Op) {
	computeI32(cpu, op, func(a, b int32) int32 {
		return int32(uint32(a) >> uint32(b))
	})
}
