package instr

import "github.com/temnok/rv/state"

func srl(cpu *state.CPU, op Op) {
	computeR(cpu, op, func(a, b int) int {
		b &= cpu.Xmask()

		return int(cpu.Xuint(a) >> cpu.Xuint(b))
	})
}
