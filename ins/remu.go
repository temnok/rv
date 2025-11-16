package ins

import (
	"github.com/temnok/rv/state"
)

func remu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if b == 0 {
			return a
		}

		return int(cpu.Uint(a) % cpu.Uint(b))
	})
}
