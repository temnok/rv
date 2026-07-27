package tlb

import "github.com/temnok/rv/state"

func Flush(cpu *state.CPU, leaveGlobals bool) {
	cpu.Update.Targets |= state.UpdateTLB
	cpu.Update.TLB[0].Vas = 0
}
