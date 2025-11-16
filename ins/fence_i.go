package ins

import "github.com/temnok/rv/state"

func fence_i(cpu *state.CPU, op Op) {
	cpu.Update.ICache.Clear()
}
