package inst

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func illegal(cpu *state.CPU, op Op) {
	trap.Enter(cpu, trap.IllegalIstruction, 0)
}

func illegal1(cpu *state.CPU) {
	trap.Enter(cpu, trap.IllegalIstruction, 0)
}
