package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func sfence_vma(cpu *state.CPU, op Op) {
	cpu.TLB.Flush()
	cpu.Update.ICache.Clear()

	if cpu.Priv == state.PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTVM) == 1 {
		illegal(cpu, op)
	}
}
