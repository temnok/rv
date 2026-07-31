package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func execLoad(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rd, rs1, offset := op.rd(), op.rs1(), imm.I(op.code())

	switch op.f3() {
	case 0:
		ctx.LB(rd, rs1, offset)
	case 1:
		ctx.LH(rd, rs1, offset)
	case 2:
		ctx.LW(rd, rs1, offset)
	case 3:
		ctx.LD(rd, rs1, offset)
	case 4:
		ctx.LBU(rd, rs1, offset)
	case 5:
		ctx.LHU(rd, rs1, offset)
	case 6:
		ctx.LWU(rd, rs1, offset)
	}
}

func memRead(ctx *context, address, width int) (int, bool) {
	val := mem.Read((*state.CPU)(ctx), address, width)

	if trap.IsEntered((*state.CPU)(ctx)) {
		return 0, false
	}

	return val, true
}
