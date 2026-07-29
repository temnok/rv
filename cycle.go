package rv

import (
	"github.com/temnok/rv/device"
	"github.com/temnok/rv/inst"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Cycle(cpu *state.CPU) int {
	cpu.Update.Targets = 0
	opcode := -1

	if trap.CheckPendingInterrupts(cpu); !trap.IsEntered(cpu) {

		if opcode = mem.Fetch(cpu, cpu.PC); !trap.IsEntered(cpu) {

			inst.Exec(cpu, opcode)

		}

	}

	device.Sync(cpu)

	state.Update(cpu)

	return opcode
}
