package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execStore(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rs2, rs1, offset := op.rs2(), op.rs1(), imm.S(op.code())

	switch op.f3() {
	case 0:
		ctx.SB(rs2, rs1, offset)
	case 1:
		ctx.SH(rs2, rs1, offset)
	case 2:
		ctx.SW(rs2, rs1, offset)
	case 3:
		ctx.SD(rs2, rs1, offset)
	}
}
