package instr

import (
	"github.com/temnok/rv/state"
)

func Lh(cpu *state.CPU, op Op) {
	load(cpu, op, 2, func(val int) int {
		return int(int16(val))
	})
}
