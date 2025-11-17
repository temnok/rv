package ins

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ecall(cpu *state.CPU, op Op) {
	trap.EnterWithoutTval(cpu, trap.EnvironmentCallFromUMode+cpu.Priv)
}
