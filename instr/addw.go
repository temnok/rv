package instr

import "github.com/temnok/rv/state"

func Addw(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		return int(int32(a) + int32(b))
	})
}
