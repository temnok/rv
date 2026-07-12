package ins

import (
	"github.com/temnok/rv/state"
)

func srl(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= 63

		return int(uint(a) >> uint(b))
	})
}
