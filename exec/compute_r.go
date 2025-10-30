package exec

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeR(cpu *state.State, op instr.Op) {
	f7, rd, rs1, rs2 := op.F7(), op.Rd(), op.Rs1(), op.Rs2()

	if f7 == 1 {
		ComputeM(cpu, rs2, rs1, op.F3(), rd)
		return
	}

	if f7&^0b0100000 != 0 {
		trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		return
	}

	switch bi.T(f7, 5)<<3 | op.F3() {
	case 0b_000:
		instr.Add(cpu, rd, rs1, rs2)
	case 0b_1_000:
		instr.Sub(cpu, rd, rs1, rs2)
	case 0b_001:
		instr.Sll(cpu, rd, rs1, rs2)
	case 0b_010:
		instr.Slt(cpu, rd, rs1, rs2)
	case 0b_011:
		instr.Sltu(cpu, rd, rs1, rs2)
	case 0b_100:
		instr.Xor(cpu, rd, rs1, rs2)
	case 0b_101:
		instr.Srl(cpu, rd, rs1, rs2)
	case 0b_1_101:
		instr.Sra(cpu, rd, rs1, rs2)
	case 0b_110:
		instr.Or(cpu, rd, rs1, rs2)
	case 0b_111:
		instr.And(cpu, rd, rs1, rs2)
	default:
		trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
	}
}
