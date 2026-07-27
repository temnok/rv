package mem

import (
	"github.com/temnok/rv/dev"
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Read(cpu *state.CPU, va int, width int) int {
	pa, _ := translateSv39(cpu, va, accessLoad)
	if trap.IsEntered(cpu) {
		return 0
	}

	if va&(width-1) != 0 {
		trap.Enter(cpu, trap.LoadAddressMisaligned, va)
		return 0
	}

	if pa < ram.BaseAddr {
		return dev.Read(cpu, pa)
	} else {
		return ram.Read(cpu, pa, width)
	}
}
