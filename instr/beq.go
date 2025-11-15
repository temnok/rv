package instr

import (
	"github.com/temnok/rv/state"
)

func beq(cpu *state.CPU, op Op) {
	branch(cpu, op, func(a, b int) bool {
		return a == b
	})
}
