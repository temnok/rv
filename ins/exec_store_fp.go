package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func execStoreFP(cpu *state.CPU, op Op) {
	if fpDisabled := cpu.CSR.Mstatus>>csr.MstatusFS&3 == 0; fpDisabled {
		illegal(cpu, op)
		return
	}

	ins := illegal

	switch op.f3() {
	case 2:
		ins = fsw

	case 3:
		ins = fsd
	}

	ins(cpu, op)
}
