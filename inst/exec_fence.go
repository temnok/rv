package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func execFence(cpu *state.CPU, op Op) {
	imm, rd, rs1 := imm.I(op.code()), op.rd(), op.rs1()

	switch op.f3() {
	case 0b_000:
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			illegal(cpu, op)
			return
		}

		fence(cpu, op)

	case 0b_001:
		if imm != 0 || rs1 != 0 || rd != 0 {
			illegal(cpu, op)
			return
		}

		fence_i(cpu, op)

	default:
		illegal(cpu, op)
	}
}
