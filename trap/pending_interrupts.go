package trap

import (
	"github.com/temnok/rv/arch"
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func OnPendingInterrupts(cpu *state.CPU) {
	cpu.Bus.NotifyInterrupts()

	mi := cpu.CSR.Mip & cpu.CSR.Mie

	if mi == 0 {
		return
	}

	for i := 9; i >= 3; i -= 2 {
		if bi.T(mi, i) == 0 {
			continue
		}

		priv := state.PrivM
		if bi.T(cpu.CSR.Mideleg, i) == 1 {
			priv = state.PrivS
		}

		mcauseI := arch.XMask
		if (priv == cpu.Priv && bi.T(cpu.CSR.Mstatus, priv) == 1) || priv > cpu.Priv {
			EnterWithoutTval(cpu, -1<<mcauseI|i)

			return
		}
	}
}
