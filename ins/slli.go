package ins

import (
	"github.com/temnok/rv/state"
)

func slli(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		return a << b
	})
}
