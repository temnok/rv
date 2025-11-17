package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func csrrs(cpu *state.CPU, op Op) {
	r, set := csrRS(cpu, op)
	old := 0

	if !csr.Read(cpu, r, &old) {
		illegal(cpu, op)
		return
	}

	if set != 0 {
		if !csr.Write(cpu, r, old|set) {
			illegal(cpu, op)
			return
		}
	}

	cpu.Xset(op.rd(), old)
}
