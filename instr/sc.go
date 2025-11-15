package instr

import (
	"github.com/temnok/rv/state"
)

func sc(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		if !cpu.Reserved || cpu.ReservedAddr != addr {
			*old = 1
			return false
		}

		cpu.Update.Reserved = false

		*old = 0
		return true
	})
}
