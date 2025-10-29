package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000:
		instr.Addi(&cpu.State, rd, rs1, imm)
	case 0b_001:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Slli(&cpu.State, rd, rs1, imm)
		default:
			trap.EnterWithoutTval(&cpu.State, ExceptionIllegalIstruction)
		}
	case 0b_010:
		instr.Slti(&cpu.State, rd, rs1, imm)
	case 0b_011:
		instr.Sltiu(&cpu.State, rd, rs1, imm)
	case 0b_100:
		instr.Xori(&cpu.State, rd, rs1, imm)
	case 0b_101:
		switch imm &^ cpu.Xmask() {
		case 0:
			instr.Srli(&cpu.State, rd, rs1, imm)
		case 0b_010000000000:
			instr.Srai(&cpu.State, rd, rs1, imm&cpu.Xmask())
		default:
			trap.EnterWithoutTval(&cpu.State, ExceptionIllegalIstruction)
		}
	case 0b_110:
		instr.Ori(&cpu.State, rd, rs1, imm)
	case 0b_111:
		instr.Andi(&cpu.State, rd, rs1, imm)
	}
}

func (cpu *CPU) execComputeR(f7, rs2, rs1, f3, rd int) {
	if f7 == 1 {
		cpu.execComputeM(rs2, rs1, f3, rd)
		return
	}

	op := bi.T(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		trap.EnterWithoutTval(&cpu.State, ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000:
		instr.Add(&cpu.State, rd, rs1, rs2)
	case 0b_1_000:
		instr.Sub(&cpu.State, rd, rs1, rs2)
	case 0b_001:
		instr.Sll(&cpu.State, rd, rs1, rs2)
	case 0b_010:
		instr.Slt(&cpu.State, rd, rs1, rs2)
	case 0b_011:
		instr.Sltu(&cpu.State, rd, rs1, rs2)
	case 0b_100:
		instr.Xor(&cpu.State, rd, rs1, rs2)
	case 0b_101:
		instr.Srl(&cpu.State, rd, rs1, rs2)
	case 0b_1_101:
		instr.Sra(&cpu.State, rd, rs1, rs2)
	case 0b_110:
		instr.Or(&cpu.State, rd, rs1, rs2)
	case 0b_111:
		instr.And(&cpu.State, rd, rs1, rs2)
	default:
		trap.EnterWithoutTval(&cpu.State, ExceptionIllegalIstruction)
	}
}
