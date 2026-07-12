package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Write(cpu *state.CPU, virtAddr, width, data int) {
	var physAddr int
	if translate.Sv(cpu, virtAddr, &physAddr, state.AccessStore); trap.IsEntered(cpu) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.StoreAMOAddressMisaligned, virtAddr)
		return
	}

	if device := cpu.DeviceAtAddress(physAddr); device != nil {
		device(physAddr, width, true, data)
		return
	}

	trap.Enter(cpu, trap.StoreAMOAccessFault, virtAddr)
}
