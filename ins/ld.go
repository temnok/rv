package ins

import (
	"github.com/temnok/rv/state"
)

func ld(cpu *state.CPU, op Op) {
	load(cpu, op, 8, func(val int) int {
		return val
	})
}
