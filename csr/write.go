package csr

import (
	"github.com/temnok/rv/state"
)

func Write(cpu *state.State, i, val int) bool {
	reg, mask, shift := addr(cpu, i, true)

	if reg == nil {
		return false
	}

	cpu.Update.CRegPtr = reg
	cpu.Update.CReg = i
	cpu.Update.CVal = *reg&^mask | (val<<shift)&mask

	return true
}
