package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func csrrw(cpu *state.CPU, op Op) {
	reg, set := csrRS(cpu, op)
	old := 0

	if op.rd() != 0 {
		if !csr.Read(cpu, reg, &old) {
			illegal(cpu, op)
			return
		}
	}

	if !csr.Write(cpu, reg, set) {
		illegal(cpu, op)
		return
	}

	cpu.Xset(op.rd(), old)
}
