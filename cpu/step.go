package cpu

import (
	"github.com/temnok/rv/dev"
	"github.com/temnok/rv/state"
)

func Step(cpu *state.CPU) int {
	cpu.Update.Targets = 0

	opcode := Evaluate(cpu)

	dev.IncrementTimer(cpu)
	dev.ProcessIO(cpu)

	Update(cpu)

	return opcode
}
