package instr

import (
	"github.com/temnok/rv/state"
)

func Lhu(cpu *state.CPU, op Op) {
	load(cpu, op, 2, func(val int) int {
		return int(uint16(val))
	})
}
