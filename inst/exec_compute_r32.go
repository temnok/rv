package inst

import (
	"github.com/temnok/rv/state"
)

func execComputeR32(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rd, rs1, rs2 := op.rd(), op.rs1(), op.rs2()

	switch f3 := op.f3(); op.f7() {
	case 0:
		switch f3 {
		case 0:
			ctx.ADDW(rd, rs1, rs2)
		case 1:
			ctx.SLLW(rd, rs1, rs2)
		case 5:
			ctx.SRLW(rd, rs1, rs2)
		}

	case 1:
		switch f3 {
		case 0:
			ctx.MULW(rd, rs1, rs2)
		case 4:
			ctx.DIVW(rd, rs1, rs2)
		case 5:
			ctx.DIVUW(rd, rs1, rs2)
		case 6:
			ctx.REMW(rd, rs1, rs2)
		case 7:
			ctx.REMUW(rd, rs1, rs2)
		}

	case 0x20:
		switch f3 {
		case 0:
			ctx.SUBW(rd, rs1, rs2)
		case 5:
			ctx.SRAW(rd, rs1, rs2)
		}
	}
}
