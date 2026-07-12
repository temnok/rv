package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Read(cpu *state.CPU, virtAddr int, width int) int {
	var physAddr int
	if physAddr = translate.Sv(cpu, virtAddr, state.AccessLoad); trap.IsEntered(cpu) {
		return 0
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.LoadAddressMisaligned, virtAddr)
		return 0
	}

	device := cpu.DeviceAtAddress(physAddr)

	if device == nil {
		trap.Enter(cpu, trap.LoadAccessFault, virtAddr)
		return 0
	}

	return device(physAddr, width, false, 0)
}
