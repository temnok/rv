package exec

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeR64(cpu *state.CPU, op instr.Op) {
	if !cpu.LenIs64() {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	f7, f3 := op.F7(), op.F3()

	if f7 == 1 {
		ComputeM64(cpu, op)
		return
	}

	if f7&^0b0100000 != 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch bi.T(f7, 5)<<3 | f3 {
	case 0b_000:
		instr.Addw(cpu, op)
	case 0b_1_000:
		instr.Subw(cpu, op)
	case 0b_001:
		instr.Sllw(cpu, op)
	case 0b_101:
		instr.Srlw(cpu, op)
	case 0b_1_101:
		instr.Sraw(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
