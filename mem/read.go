package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Read(cpu *state.State, virtAddr int, data *int, width int) {
	var physAddr int
	if translate.Sv(cpu, virtAddr, &physAddr, state.AccessRead); trap.IsEntered(cpu) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.LoadAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Read(physAddr, data, width) {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
	}
}
