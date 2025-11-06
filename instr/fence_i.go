package instr

import "github.com/temnok/rv/state"

func Fence_I(cpu *state.CPU, op Op) {
	cpu.Update.ICache.Clear()
}
