package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Branch(cpu *state.State, op instr.Op) {
	switch op.F3() {
	case 0b_000:
		instr.Beq(cpu, op)
	case 0b_001:
		instr.Bne(cpu, op)
	case 0b_100:
		instr.Blt(cpu, op)
	case 0b_101:
		instr.Bge(cpu, op)
	case 0b_110:
		instr.Bltu(cpu, op)
	case 0b_111:
		instr.Bgeu(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
