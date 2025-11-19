package ins

import (
	"github.com/temnok/rv/state"
)

func csrrs(cpu *state.CPU, op Op) {
	csrAccess(cpu, op, true, false, func(set, old int) int {
		return old | set
	})
}
