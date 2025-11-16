package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func Csrrc(cpu *state.CPU, op Op) {
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
