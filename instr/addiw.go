package instr

import (
	"github.com/temnok/rv/state"
)

func Addiw(cpu *state.CPU, op Op) {
	computeI32(cpu, op, func(a, b int32) int32 {
		return a + b
	})
}
