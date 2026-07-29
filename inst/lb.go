package inst

import (
	"github.com/temnok/rv/state"
)

func lb(cpu *state.CPU, op Op) {
	load(cpu, op, 1, func(val int) int {
		return int(int8(val))
	})
}
