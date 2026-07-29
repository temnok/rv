package inst

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func sfence_vma(cpu *state.CPU, op Op) {
	if cpu.CSR.Priv == csr.PrivS && cpu.CSR.Mstatus>>csr.MstatusTVM&1 == 1 {
		illegal(cpu, op)
	}
}
