package instr

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func ComputeR(cpu *state.CPU, op Op) {
	f7 := op.F7()

	switch {
	case f7 == 1:
		computeM(cpu, op)
		return

	case f7&^0b0100000 != 0:
		illegal(cpu, op)
		return
	}

	switch bi.T(f7, 5)<<3 | op.F3() {
	case 0b_000:
		add(cpu, op)
	case 0b_1_000:
		sub(cpu, op)
	case 0b_001:
		sll(cpu, op)
	case 0b_010:
		slt(cpu, op)
	case 0b_011:
		sltu(cpu, op)
	case 0b_100:
		xor(cpu, op)
	case 0b_101:
		srl(cpu, op)
	case 0b_1_101:
		sra(cpu, op)
	case 0b_110:
		or(cpu, op)
	case 0b_111:
		and(cpu, op)
	default:
		illegal(cpu, op)
	}
}
