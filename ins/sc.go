package ins

import (
	"github.com/temnok/rv/state"
)

func sc(cpu *state.CPU, op Op) {
	atomic(cpu, op, false, func(cpu *state.CPU, addr int, val, old *int) bool {
		if cpu.Reservation == -1 || cpu.Reservation != addr {
			*old = 1
			return false
		}

		cpu.Update.Targets |= state.UpdateReservation
		cpu.Update.Reservation = -1

		*old = 0
		return true
	})
}
