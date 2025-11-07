package instr

import (
	"github.com/temnok/rv/state"
)

func Sh(cpu *state.CPU, op Op) {
	store(cpu, op, 2)
}
