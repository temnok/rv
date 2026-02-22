package ins

import (
	"github.com/temnok/rv/arch"
	"github.com/temnok/rv/state"
)

func srl(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= arch.XMask

		return int(uint(a) >> uint(b))
	})
}
