package instr

import (
	"github.com/temnok/rv/state"
)

func Lwu(cpu *state.CPU, op Op) {
	load(cpu, op, 4, func(val int) int {
		return int(uint32(val))
	})
}
