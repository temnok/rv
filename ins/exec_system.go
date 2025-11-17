package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var systemIns = []func(*state.CPU, Op){
	0: systemSpecial,
	1: csrrw,
	2: csrrs,
	3: csrrc,
	4: illegal,
	5: csrrw,
	6: csrrs,
	7: csrrc,
}

func execSystem(cpu *state.CPU, op Op) {
	systemIns[op.f3()](cpu, op)
}

func systemSpecial(cpu *state.CPU, op Op) {
	imm, rd := imm.I(op.code()), op.rd()

	if rd != 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000:
		ecall(cpu, op)

	case 0b_0000_000_00001:
		ebreak(cpu, op)

	case 0b_0001_000_00010:
		sret(cpu, op)

	case 0b_0001_000_00101:
		wfi(cpu, op)

	case 0b_0011_000_00010:
		mret(cpu, op)

	default:
		switch bi.Ts(imm, 5, 7) {
		case 0b_0001_001:
			sfence_vma(cpu, op)

		default:
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		}
	}
}
