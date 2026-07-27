package trap

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func Exit(cpu *state.CPU, retPriv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapped := cpu.CSR.Priv == csr.PrivS && bit.IsSet(cpu.CSR.Mstatus, csr.MstatusTSR)

	if trapped || retPriv > cpu.CSR.Priv {
		Enter(cpu, IllegalIstruction, 0)
		return
	}

	ms := cpu.CSR.Mstatus

	switch retPriv {
	case csr.PrivM:
		cpu.Update.PC = cpu.CSR.Mepc
		cpu.Update.Priv = bit.GetN(cpu.CSR.Mstatus, csr.MstatusMPP, 2)

		ms &^= 3 << csr.MstatusMPP
		ms |= 1<<csr.MstatusMPIE | (ms>>csr.MstatusMPIE&1)<<csr.MstatusMIE

	case csr.PrivS:
		cpu.Update.PC = cpu.CSR.Sepc
		cpu.Update.Priv = bit.GetN(cpu.CSR.Mstatus, csr.MstatusSPP, 1)

		ms &^= 1 << csr.MstatusSPP
		ms |= 1<<csr.MstatusSPIE | (ms>>csr.MstatusSPIE&1)<<csr.MstatusSIE
	}

	if cpu.CSR.Priv != csr.PrivM {
		ms &^= 1 << csr.MstatusMPRV
	}

	cpu.Update.Targets |= state.UpdatePriv | state.UpdateMstatus
	cpu.Update.Mstatus = ms
}
