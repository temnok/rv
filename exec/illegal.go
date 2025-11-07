package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func instrIllegal(cpu *state.CPU, op instr.Op) {
	trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
}
