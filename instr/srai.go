package instr

import (
	"github.com/temnok/rv/state"
)

func srai(cpu *state.CPU, op Op) {
	computeI(cpu, op, func(a, b int) int {
		b &= cpu.Xmask()

		return a >> b
	})
}
