package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
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
		illegal(cpu, op)
		return
	}

	ins := illegal

	switch imm {
	case 0:
		ins = ecall
	case 1:
		ins = ebreak
	case 1<<8 | 2:
		ins = sret
	case 1<<8 | 5:
		ins = wfi
	case 3<<8 | 2:
		ins = mret
	default:
		if bi.Ts(imm, 5, 7) == 9 {
			ins = sfence_vma
		}
	}

	ins(cpu, op)
}

func csrAccess(cpu *state.CPU, op Op, enforceRead, enforceWrite bool, f func(set, old int) int) {
	reg, set, old := bi.Ts(imm.I(op.code()), 0, 12), op.rs1(), 0

	if op.f3()&4 == 0 {
		set = cpu.X[set]
	}

	if op.rd() != 0 || enforceRead {
		if !csr.Read(cpu, reg, &old) {
			illegal(cpu, op)
			return
		}
	}

	if set != 0 || enforceWrite {
		if !csr.Write(cpu, reg, f(set, old)) {
			illegal(cpu, op)
			return
		}
	}

	cpu.Xset(op.rd(), old)
}
