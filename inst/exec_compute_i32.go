package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execComputeI32(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)
	rd, rs1, imm := op.rd(), op.rs1(), imm.I(op.code())

	switch op.f3() {
	case 0:
		ctx.ADDIW(rd, rs1, imm)
	case 1:
		ctx.SLLIW(rd, rs1, imm)
	case 5:
		switch imm &^ 0x1F {
		case 0:
			ctx.SRLIW(rd, rs1, imm)
		case 0x400:
			ctx.SRAIW(rd, rs1, imm)
		}
	}
}
