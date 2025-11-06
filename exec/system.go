package exec

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func System(cpu *state.CPU, op instr.Op) {
	if op.F3() == 0 {
		systemSpecial(cpu, op)
	} else {
		systemCSR(cpu, op)
	}
}

func systemSpecial(cpu *state.CPU, op instr.Op) {
	imm, rd := imm.I(op.Code()), op.Rd()

	if rd != 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000:
		instr.Ecall(cpu, op)

	case 0b_0000_000_00001:
		instr.Ebreak(cpu, op)

	case 0b_0001_000_00010:
		instr.Sret(cpu, op)

	case 0b_0001_000_00101:
		instr.Wfi(cpu, op)

	case 0b_0011_000_00010:
		instr.Mret(cpu, op)

	default:
		switch bi.Ts(imm, 5, 7) {
		case 0b_0001_001:
			instr.Sfence_vma(cpu, op)

		default:
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		}
	}
}

func systemCSR(cpu *state.CPU, op instr.Op) {
	switch op.F3() & 3 {
	case 1:
		instr.Csrrw(cpu, op)
	case 2:
		instr.Csrrs(cpu, op)
	case 3:
		instr.Csrrc(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
