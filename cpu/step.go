package cpu

import (
	"github.com/temnok/rv/ins"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Step(cpu *state.CPU) int {
	cpu.Update.Targets = state.UpdatePC

	opcode := FetchAndExec(cpu)

	UpdateState(cpu)

	return opcode
}

func FetchAndExec(cpu *state.CPU) int {
	if trap.OnPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return -1
	}

	var opcode int
	if opcode = mem.Fetch(cpu, cpu.PC); trap.IsEntered(cpu) {
		return -1
	}

	ins.Exec(cpu, opcode)

	return opcode
}
