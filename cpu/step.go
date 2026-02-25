package cpu

import (
	"github.com/temnok/rv/ins"
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

	if cpu.InstrCount%20 == 0 {
		if trap.OnPendingInterrupts(cpu); trap.IsEntered(cpu) {
			return 0
		}
	}

	var opcode int
	if mem.Fetch(cpu, cpu.PC, &opcode); trap.IsEntered(cpu) {
		return 0
	}

	ins.Exec(cpu, opcode)

	return opcode
}
