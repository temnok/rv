package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execComputeI(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rd, rs1, imm := op.rd(), op.rs1(), imm.I(op.code())

	switch op.f3() {
	case 0:
		ctx.ADDI(rd, rs1, imm)
	case 1:
		ctx.SLLI(rd, rs1, imm)
	case 2:
		ctx.SLTI(rd, rs1, imm)
	case 3:
		ctx.SLTIU(rd, rs1, imm)
	case 4:
		ctx.XORI(rd, rs1, imm)
	case 5:
		switch imm &^ 0x3F {
		case 0:
			ctx.SRLI(rd, rs1, imm)
		case 0x400:
			ctx.SRAI(rd, rs1, imm)
		}
	case 6:
		ctx.ORI(rd, rs1, imm)
	case 7:
		ctx.ANDI(rd, rs1, imm)
	}
}
