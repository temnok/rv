package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Read(cpu *state.CPU, virtAddr int, data *int, width int) {
	var physAddr int
	if translate.Sv(cpu, virtAddr, &physAddr, state.AccessRead); trap.IsEntered(cpu) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.LoadAddressMisaligned, virtAddr)
		return
	}

	var ok bool
	if *data, ok = cpu.Bus.Read(physAddr, width); !ok {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
	}
}
