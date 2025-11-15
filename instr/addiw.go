package instr

import (
	"github.com/temnok/rv/state"
)

func Addiw(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		return int(int32(a) + int32(b))
	})
}
