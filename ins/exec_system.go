package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func execSystem(cpu *state.CPU, op Op) {
	if op.F3() == 0 {
		systemSpecial(cpu, op)
	} else {
		systemCSR(cpu, op)
	}
}

func systemSpecial(cpu *state.CPU, op Op) {
	imm, rd := imm.I(op.Code()), op.Rd()

	if rd != 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000:
		Ecall(cpu, op)

	case 0b_0000_000_00001:
		Ebreak(cpu, op)

	case 0b_0001_000_00010:
		Sret(cpu, op)

	case 0b_0001_000_00101:
		wfi(cpu, op)

	case 0b_0011_000_00010:
		Mret(cpu, op)

	default:
		switch bi.Ts(imm, 5, 7) {
		case 0b_0001_001:
			Sfence_vma(cpu, op)

		default:
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		}
	}
}

func systemCSR(cpu *state.CPU, op Op) {
	switch op.F3() & 3 {
	case 1:
		Csrrw(cpu, op)
	case 2:
		Csrrs(cpu, op)
	case 3:
		Csrrc(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
