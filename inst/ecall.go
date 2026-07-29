package inst

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ecall(cpu *state.CPU, op Op) {
	trap.Enter(cpu, trap.EnvironmentCallFromUMode+cpu.CSR.Priv, 0)
}
