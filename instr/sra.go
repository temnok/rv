package instr

import "github.com/temnok/rv/state"

func sra(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= cpu.Xmask()

		return a >> b
	})
}
