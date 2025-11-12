package instr

import (
	"github.com/temnok/rv/state"
)

func Amomaxu(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		if cpu.Xuint(*old) > cpu.Xuint(*val) {
			*val = *old
		}

		return true
	})
}
