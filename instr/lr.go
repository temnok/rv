package instr

import (
	"github.com/temnok/rv/state"
)

func Lr(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		cpu.Update.Reserved = true
		cpu.Update.ReservedAddr = addr

		return false
	})
}
