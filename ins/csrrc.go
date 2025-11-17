package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func csrrc(cpu *state.CPU, op Op) {
	r, s := csrRS(cpu, op)
	val := 0

	if !csr.Read(cpu, r, &val) {
		illegal(cpu, op)
		return
	}

	if s != 0 {
		if !csr.Write(cpu, r, val&^s) {
			illegal(cpu, op)
			return
		}
	}

	cpu.Xset(op.Rd(), val)
}

func csrRS(cpu *state.CPU, op Op) (int, int) {
	r := bi.Ts(imm.I(op.Code()), 0, 12)

	s := op.Rs1()
	if (op.F3() & 4) == 0 {
		s = cpu.X[s]
	}

	return r, s
}
