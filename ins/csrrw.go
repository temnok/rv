package ins

import (
	"github.com/temnok/rv/state"
)

func csrrw(cpu *state.CPU, op Op) {
	csrAccess(cpu, op, false, true, func(set, old int) int {
		return set
	})
}
