package inst

import (
	"github.com/temnok/rv/state"
)

func lwu(cpu *state.CPU, op Op) {
	load(cpu, op, 4, func(val int) int {
		return int(uint32(val))
	})
}
