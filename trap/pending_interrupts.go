package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:intr_M-mode_pri
var interruptPriorityOrder = []int{
	csr.MipMEIP,
	csr.MipMSIP,
	csr.MipMTIP,
	csr.MipSEIP,
	csr.MipSSIP,
	csr.MipSTIP,
}

// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:mcause_exccode_enc_img
func CheckPendingInterrupts(cpu *state.CPU) {
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-9-machine-interrupt-mip-and-mie-registers
	mi := cpu.CSR.Mip & cpu.CSR.Mie

	if mi == 0 {
		return
	}

	for _, i := range interruptPriorityOrder {
		if mi>>i&1 == 0 {
			continue
		}

		priv := csr.PrivM
		xIE := csr.MstatusMIE
		if delegateToSmode := cpu.CSR.Mideleg>>i&1 == 1; delegateToSmode {
			priv = csr.PrivS
			xIE = csr.MstatusSIE
		}

		// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#privstack
		if priv > cpu.CSR.Priv || priv == cpu.CSR.Priv && cpu.CSR.Mstatus>>xIE&1 == 1 {
			Enter(cpu, -1<<csr.McauseI|i, 0)
			break
		}
	}
}
