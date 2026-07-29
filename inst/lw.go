package inst

import (
	"github.com/temnok/rv/state"
)

func lw(cpu *state.CPU, op Op) {
	load(cpu, op, 4, func(val int) int {
		return int(int32(val))
	})
}
