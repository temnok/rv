package instr

import (
	"github.com/temnok/rv/state"
)

func Sd(cpu *state.CPU, op Op) {
	store(cpu, op, 8)
}
