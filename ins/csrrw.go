package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func csrrw(cpu *state.CPU, op Op) {
	r, s := csrRS(cpu, op)
	val := 0

	if op.Rd() != 0 {
		if !csr.Read(cpu, r, &val) {
			illegal(cpu, op)
			return
		}
	}

	if !csr.Write(cpu, r, s) {
		illegal(cpu, op)
		return
	}

	cpu.Xset(op.Rd(), val)
}
