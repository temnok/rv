package inst

import (
	"github.com/temnok/rv/state"
)

// https://docs.riscv.org/reference/isa/v20260120/unpriv/rv32.html#fence
func execFence(cpu *state.CPU, op Op) {
	ctx := (*context)(cpu)

	switch op.f3() {
	case 0b_000:

		switch mode := op.code() >> 28 & 0xF; mode {
		case 0:
			pred, succ := op.code()>>24&0xF, op.code()>>20&0xF
			ctx.FENCE(pred, succ)
		case 8:
			ctx.FENCE_TSO()
		}

	case 0b_001:
		ctx.FENCE_I()
	}
}
