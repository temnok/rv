package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeI(cpu *state.State, op instr.Op) {
	imm := imm.I(op.Code())

	switch op.F3() {
	case 0b_000:
		instr.Addi(cpu, op)
	case 0b_001:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Slli(cpu, op)
		default:
			trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		}
	case 0b_010:
		instr.Slti(cpu, op)
	case 0b_011:
		instr.Sltiu(cpu, op)
	case 0b_100:
		instr.Xori(cpu, op)
	case 0b_101:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Srli(cpu, op)
		case 0b_010000000000:
			instr.Srai(cpu, op)
		default:
			trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		}
	case 0b_110:
		instr.Ori(cpu, op)
	case 0b_111:
		instr.Andi(cpu, op)
	}
}
