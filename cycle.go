package rv

import (
	"github.com/temnok/rv/dev"
	"github.com/temnok/rv/ins"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Cycle(cpu *state.CPU) int {
	cpu.Update.Targets = 0
	opcode := -1

	if trap.CheckPendingInterrupts(cpu); !trap.IsEntered(cpu) {

		if opcode = mem.Fetch(cpu, cpu.PC); !trap.IsEntered(cpu) {

			ins.Exec(cpu, opcode)

		}

	}

	dev.Process(cpu)

	state.Update(cpu)

	return opcode
}
