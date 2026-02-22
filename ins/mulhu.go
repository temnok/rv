package ins

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func mulhu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		hi, _ := bits.Mul64(uint64(a), uint64(b))
		return int(hi)
	})
}
