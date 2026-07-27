package tlb

import "github.com/temnok/rv/state"

func Append(cpu *state.CPU, virtAddr, shift, pte int) {
	cpu.Update.Targets |= state.UpdateTLB

	for i := len(cpu.TLB) - 1; i > 0; i-- {
		cpu.Update.TLB[i] = cpu.TLB[i-1]
	}

	cpu.Update.TLB[0].Vas = virtAddr&(-1<<shift) | shift
	cpu.Update.TLB[0].Pte = pte
}
