package rv

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execComputeM64(op instr.Op) {
	switch op.F3() {
	case 0b_000:
		instr.Mulw(cpu.State, op)
	case 0b_100:
		instr.Divw(cpu.State, op)
	case 0b_101:
		instr.Divuw(cpu.State, op)
	case 0b_110:
		instr.Remw(cpu.State, op)
	case 0b_111:
		instr.Remuw(cpu.State, op)
	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
