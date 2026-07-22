package trap

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func Exit(cpu *state.CPU, retPriv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapped := cpu.Priv == state.PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTSR) == 1

	if trapped || retPriv > cpu.Priv {
		EnterWithoutTval(cpu, IllegalIstruction)
		return
	}

	ms := cpu.CSR.Mstatus

	switch retPriv {
	case state.PrivM:
		cpu.Update.PC = cpu.CSR.Mepc
		cpu.Update.Priv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusMPP, 2)

		ms &^= 3 << csr.MstatusMPP
		ms |= 1<<csr.MstatusMPIE | (ms>>csr.MstatusMPIE&1)<<csr.MstatusMIE

	case state.PrivS:
		cpu.Update.PC = cpu.CSR.Sepc
		cpu.Update.Priv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusSPP, 1)

		ms &^= 1 << csr.MstatusSPP
		ms |= 1<<csr.MstatusSPIE | (ms>>csr.MstatusSPIE&1)<<csr.MstatusSIE
	}

	if cpu.Priv != state.PrivM {
		ms &^= 1 << csr.MstatusMPRV
	}

	cpu.Update.Targets |= state.UpdatePriv | state.UpdateMstatus
	cpu.Update.Mstatus = ms
}
