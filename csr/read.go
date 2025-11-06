package csr

import (
	"github.com/temnok/rv/state"
)

func Read(cpu *state.CPU, i int, val *int) bool {
	reg, mask, shift := addr(cpu, i, false)

	if reg == nil {
		return false
	}

	*val = (*reg & mask) >> shift

	return true
}
