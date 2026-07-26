package mem

import (
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Write(cpu *state.CPU, virtAddr, width, val int) {
	var physAddr int
	if physAddr = translate.Sv(cpu, virtAddr, state.AccessStore); trap.IsEntered(cpu) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.StoreAMOAddressMisaligned, virtAddr)
		return
	}

	if physAddr < ram.BaseAddr {
		cpu.Devices(physAddr, width, true, val)
	} else {
		ram.Write(cpu.RAM, physAddr, width, val)
	}
}
