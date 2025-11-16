package ins

import (
	"github.com/temnok/rv/state"
)

func srai(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		b &= cpu.Mask()

		return a >> b
	})
}
