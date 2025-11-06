package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeM64(cpu *state.CPU, op instr.Op) {
	switch op.F3() {
	case 0b_000:
		instr.Mulw(cpu, op)
	case 0b_100:
		instr.Divw(cpu, op)
	case 0b_101:
		instr.Divuw(cpu, op)
	case 0b_110:
		instr.Remw(cpu, op)
	case 0b_111:
		instr.Remuw(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
