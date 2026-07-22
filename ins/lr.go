package ins

import (
	"github.com/temnok/rv/state"
)

func lr(cpu *state.CPU, op Op) {
	atomic(cpu, op, true, func(cpu *state.CPU, addr int, val, old *int) bool {
		cpu.Update.Targets |= state.UpdateReservation
		cpu.Update.Reservation = addr

		return false
	})
}
