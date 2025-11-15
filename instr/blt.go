package instr

import (
	"github.com/temnok/rv/state"
)

func blt(cpu *state.CPU, op Op) {
	branch(cpu, op, func(a, b int) bool {
		return a < b
	})
}
