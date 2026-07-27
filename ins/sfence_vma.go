package ins

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/tlb"
)

func sfence_vma(cpu *state.CPU, op Op) {
	tlb.Flush(cpu, op.rs2() != 0)

	if cpu.CSR.Priv == csr.PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTVM) == 1 {
		illegal(cpu, op)
	}
}
