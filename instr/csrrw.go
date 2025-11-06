package instr

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Csrrw(cpu *state.CPU, op Op) {
	r, s := csrRS(cpu, op)
	val := 0

	if op.Rd() != 0 {
		if !csr.Read(cpu, r, &val) {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}
	}

	if !csr.Write(cpu, r, s) {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	cpu.Xset(op.Rd(), val)
}
