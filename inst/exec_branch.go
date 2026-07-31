package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execBranch(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rs1, rs2, offset := op.rs1(), op.rs2(), imm.B(op.code())

	switch op.f3() {
	case 0:
		ctx.BEQ(rs1, rs2, offset)
	case 1:
		ctx.BNE(rs1, rs2, offset)
	case 4:
		ctx.BLT(rs1, rs2, offset)
	case 5:
		ctx.BGE(rs1, rs2, offset)
	case 6:
		ctx.BLTU(rs1, rs2, offset)
	case 7:
		ctx.BGEU(rs1, rs2, offset)
	}
}
