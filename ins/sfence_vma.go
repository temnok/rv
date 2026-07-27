package ins

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func sfence_vma(cpu *state.CPU, op Op) {
	if cpu.CSR.Priv == csr.PrivS && bit.IsSet(cpu.CSR.Mstatus, csr.MstatusTVM) {
		illegal(cpu, op)
	}
}
