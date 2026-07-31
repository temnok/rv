package inst

import (
	"github.com/temnok/rv/state"
)

func execComputeR(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rd, rs1, rs2 := op.rd(), op.rs1(), op.rs2()

	switch f3 := op.f3(); op.f7() {
	case 0:
		switch f3 {
		case 0:
			ctx.ADD(rd, rs1, rs2)
		case 1:
			ctx.SLL(rd, rs1, rs2)
		case 2:
			ctx.SLT(rd, rs1, rs2)
		case 3:
			ctx.SLTU(rd, rs1, rs2)
		case 4:
			ctx.XOR(rd, rs1, rs2)
		case 5:
			ctx.SRL(rd, rs1, rs2)
		case 6:
			ctx.OR(rd, rs1, rs2)
		case 7:
			ctx.AND(rd, rs1, rs2)
		}

	case 1:
		switch f3 {
		case 0:
			ctx.MUL(rd, rs1, rs2)
		case 1:
			ctx.MULH(rd, rs1, rs2)
		case 2:
			ctx.MULHSU(rd, rs1, rs2)
		case 3:
			ctx.MULHU(rd, rs1, rs2)
		case 4:
			ctx.DIV(rd, rs1, rs2)
		case 5:
			ctx.DIVU(rd, rs1, rs2)
		case 6:
			ctx.REM(rd, rs1, rs2)
		case 7:
			ctx.REMU(rd, rs1, rs2)
		}

	case 0x20:
		switch f3 {
		case 0:
			ctx.SUB(rd, rs1, rs2)
		case 5:
			ctx.SRA(rd, rs1, rs2)
		}
	}
}
