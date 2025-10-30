package exec

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeR(cpu *state.State, op instr.Op) {
	f7 := op.F7()

	switch {
	case f7 == 1:
		ComputeM(cpu, op)
		return

	case f7&^0b0100000 != 0:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch bi.T(f7, 5)<<3 | op.F3() {
	case 0b_000:
		instr.Add(cpu, op)
	case 0b_1_000:
		instr.Sub(cpu, op)
	case 0b_001:
		instr.Sll(cpu, op)
	case 0b_010:
		instr.Slt(cpu, op)
	case 0b_011:
		instr.Sltu(cpu, op)
	case 0b_100:
		instr.Xor(cpu, op)
	case 0b_101:
		instr.Srl(cpu, op)
	case 0b_1_101:
		instr.Sra(cpu, op)
	case 0b_110:
		instr.Or(cpu, op)
	case 0b_111:
		instr.And(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
