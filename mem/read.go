package mem

import (
	"github.com/temnok/rv/device"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Read(cpu *state.CPU, virtAddr int, width int) int {
	var physAddr int
	if physAddr = translateSv39(cpu, virtAddr, accessLoad); trap.IsEntered(cpu) {
		return 0
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.LoadAddressMisaligned, virtAddr)
		return 0
	}

	if physAddr < ram.BaseAddr {
		return device.Read(cpu, physAddr)
	} else {
		return ram.Read(cpu, physAddr, width)
	}
}
