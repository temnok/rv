package instr

import (
	"github.com/temnok/rv/state"
)

func Amoxor(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		*val ^= *old

		return true
	})
}
