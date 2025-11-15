package instr

import (
	"github.com/temnok/rv/state"
)

func sw(cpu *state.CPU, op Op) {
	store(cpu, op, 4)
}
