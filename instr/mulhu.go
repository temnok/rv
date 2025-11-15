package instr

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func mulhu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		if !cpu.Xlen64() {
			a = int(uint32(a))
			b = int(uint32(b))

			return (a * b) >> 32
		}

		hi, _ := bits.Mul64(uint64(a), uint64(b))
		return int(hi)
	})
}
