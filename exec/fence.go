package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Fence(cpu *state.CPU, op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	switch op.F3() {
	case 0b_000:
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		instr.Fence(cpu, op)

	case 0b_001:
		if imm != 0 || rs1 != 0 || rd != 0 {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		instr.Fence_I(cpu, op)

	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
