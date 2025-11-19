package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func execLoadFP(cpu *state.CPU, op Op) {
	if csr.FpDisabled(cpu) {
		illegal(cpu, op)
		return
	}

	ins := illegal

	switch op.f3() {
	case 2:
		ins = flw

	case 3:
		ins = fld
	}

	ins(cpu, op)
}
