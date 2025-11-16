package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func ComputeI64(cpu *state.CPU, op instr.Op) {
	imm := imm.I(op.Code())

	if cpu.LenIs64() {
		switch op.F3() {
		case 0b_000:
			instr.Addiw(cpu, op)
		case 0b_001:
			if imm < 32 {
				instr.Slliw(cpu, op)
			}
		case 0b_101:
			if imm < 32 {
				instr.Srliw(cpu, op)
			} else if imm&^0b0100000_00000 < 32 {
				instr.Sraiw(cpu, op)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
