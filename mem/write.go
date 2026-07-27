package mem

import (
	"github.com/temnok/rv/dev"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Write(cpu *state.CPU, va, width, val int) {
	pa, _ := translateSv39(cpu, va, accessStore)
	if trap.IsEntered(cpu) {
		return
	}

	if va&(width-1) != 0 {
		trap.Enter(cpu, trap.StoreAMOAddressMisaligned, va)
		return
	}

	if pa < ram.BaseAddr {
		dev.Write(cpu, pa, val)
	} else {
		ram.Write(cpu, pa, width, val)
	}
}
