package instr

import (
	"github.com/temnok/rv/state"
)

func bgeu(cpu *state.CPU, op Op) {
	branch(cpu, op, func(a, b int) bool {
		return cpu.Xuint(a) >= cpu.Xuint(b)
	})
}
