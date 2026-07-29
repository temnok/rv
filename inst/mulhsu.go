package inst

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func mulhsu(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		hi, _ := bits.Mul64(uint64(a), uint64(b))
		s := (a >> 63) & b
		return int(hi) - s
	})
}
