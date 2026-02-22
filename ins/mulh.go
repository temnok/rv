package ins

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func mulh(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		hi, _ := bits.Mul64(uint64(a), uint64(b))
		s1 := (a >> 63) & b
		s2 := (b >> 63) & a
		return int(hi) - s1 - s2
	})
}
