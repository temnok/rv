package ins

import "github.com/temnok/rv/state"

func srl(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= cpu.Mask()

		return int(cpu.Uint(a) >> cpu.Uint(b))
	})
}
