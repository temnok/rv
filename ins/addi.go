package ins

import (
	"github.com/temnok/rv/state"
)

func addi(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		return a + b
	})
}
