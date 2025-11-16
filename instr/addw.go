package instr

import "github.com/temnok/rv/state"

func Addw(cpu *state.CPU, op Op) {
	computeR64(cpu, op, func(a, b int32) int32 {
		return a + b
	})
}
