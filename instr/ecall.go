package instr

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Ecall(cpu *state.CPU, op Op) {
	trap.EnterWithoutTval(cpu, trap.EnvironmentCallFromUMode+cpu.Priv)
}
