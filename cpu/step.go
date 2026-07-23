package cpu

import (
	"github.com/temnok/rv/ins"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Step(cpu *state.CPU) int {
	cpu.Update.Targets = 0

	opcode := EvaluateState(cpu)

	UpdateState(cpu)

	return opcode
}

func EvaluateState(cpu *state.CPU) int {
	if trap.CheckPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return -1
	}

	opcode := mem.Fetch(cpu, cpu.PC)
	if trap.IsEntered(cpu) {
		return -1
	}

	ins.Exec(cpu, opcode)

	return opcode
}
