package csr

import (
	"github.com/temnok/rv/state"
)

func Write(cpu *state.CPU, reg, val int) bool {
	if reg>>10&3 == 3 {
		return false
	}

	if _, ok := Read(cpu, reg); ok {
		cpu.Cset(reg, val)

		return true
	}

	return false
}
