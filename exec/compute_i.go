package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeI(cpu *state.State, op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	switch op.F3() {
	case 0b_000:
		instr.Addi(cpu, rd, rs1, imm)
	case 0b_001:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Slli(cpu, rd, rs1, imm)
		default:
			trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		}
	case 0b_010:
		instr.Slti(cpu, rd, rs1, imm)
	case 0b_011:
		instr.Sltiu(cpu, rd, rs1, imm)
	case 0b_100:
		instr.Xori(cpu, rd, rs1, imm)
	case 0b_101:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Srli(cpu, rd, rs1, imm)
		case 0b_010000000000:
			instr.Srai(cpu, rd, rs1, imm&cpu.Xmask())
		default:
			trap.EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		}
	case 0b_110:
		instr.Ori(cpu, rd, rs1, imm)
	case 0b_111:
		instr.Andi(cpu, rd, rs1, imm)
	}
}
