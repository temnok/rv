package mem

import (
	"github.com/temnok/rv/device"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Write(cpu *state.CPU, virtAddr, width, val int) {
	var physAddr int
	if physAddr = translateSv39(cpu, virtAddr, accessStore); trap.IsEntered(cpu) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu, trap.StoreAMOAddressMisaligned, virtAddr)
		return
	}

	if physAddr < ram.BaseAddr {
		device.Write(cpu, physAddr, val)
	} else {
		ram.Write(cpu, physAddr, width, val)
	}
}
