package instr

import "github.com/temnok/rv/state"

func Fence_I(cpu *state.State, op Op) {
	cpu.Update.ICache.Clear()
}
