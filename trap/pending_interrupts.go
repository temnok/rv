package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func OnPendingInterrupts(cpu *state.CPU) {
	mi := cpu.CSR.Mip & cpu.CSR.Mie

	if mi == 0 {
		return
	}

	for i := 9; i >= 1; i -= 2 {
		if mi>>i&1 == 0 {
			continue
		}

		priv := state.PrivM
		if cpu.CSR.Mideleg>>i&1 == 1 {
			priv = state.PrivS
		}

		if priv > cpu.Priv || priv == cpu.Priv && cpu.CSR.Mstatus>>priv&1 == 1 {
			Enter(cpu, -1<<csr.McauseI|i, 0)

			break
		}
	}
}
