package inst

import (
	"github.com/temnok/rv/state"
)

func csrrc(cpu *state.CPU, op Op) {
	csrAccess(cpu, op, true, false, func(set, old int) int {
		return old &^ set
	})
}
