package cpu

import (
	"github.com/temnok/rv/ins"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Step(cpu *state.CPU) int {
	opcode := fetchAndExec(cpu)

	updateState(cpu)

	return opcode
}

func fetchAndExec(cpu *state.CPU) int {
	if trap.OnPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return -1
	}

	var opcode int
	if opcode = mem.Fetch(cpu, cpu.PC); trap.IsEntered(cpu) {
		return -1
	}

	if ins.Exec(cpu, opcode); trap.IsEntered(cpu) {
		return -1
	}

	return opcode
}
