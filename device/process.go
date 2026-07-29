package device

import "github.com/temnok/rv/state"

func Sync(cpu *state.CPU) {
	syncTimer(cpu)
	syncUartInput(cpu)
	syncUartOutput(cpu)
}
