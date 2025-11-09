package cpu

import (
	"github.com/temnok/rv/exec"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Step(cpu *state.CPU) bool {
	//return debugStep(cpu)

	InnerStep(cpu)
	return true
}

func InnerStep(cpu *state.CPU) int {
	updateState(cpu)

	if trap.OnPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return 0
	}

	var opcode int
	if mem.Fetch(cpu, cpu.PC, &opcode); trap.IsEntered(cpu) {
		return 0
	}

	exec.Exec(cpu, opcode)

	return opcode
}
