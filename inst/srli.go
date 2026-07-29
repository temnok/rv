package inst

import (
	"github.com/temnok/rv/state"
)

func srli(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		return int(uint(a) >> uint(b))
	})
}
