package instr

import (
	"github.com/temnok/rv/state"
)

func lbu(cpu *state.CPU, op Op) {
	load(cpu, op, 1, func(val int) int {
		return int(uint8(val))
	})
}
