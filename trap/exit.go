package trap

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func Exit(cpu *state.State, retPriv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trapped := cpu.Priv == state.PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTSR) == 1

	if trapped || retPriv > cpu.Priv {
		EnterWithoutTval(cpu, state.ExceptionIllegalIstruction)
		return
	}

	cpu.Update.TrapExit = true

	switch retPriv {
	case state.PrivM:
		cpu.Update.TrapPC = cpu.CSR.Mepc
		cpu.Update.TrapPriv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusMPP, 2)

		mie := bi.T(cpu.CSR.Mstatus, csr.MstatusMPIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(3<<csr.MstatusMPP) |
			(1<<csr.MstatusMPIE | mie<<csr.MstatusMIE)

	case state.PrivS:
		cpu.Update.TrapPC = cpu.CSR.Sepc
		cpu.Update.TrapPriv = bi.Ts(cpu.CSR.Mstatus, csr.MstatusSPP, 1)

		sie := bi.T(cpu.CSR.Mstatus, csr.MstatusSPIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(1<<csr.MstatusSPP) |
			(1<<csr.MstatusSPIE | sie<<csr.MstatusSIE)
	}

	if cpu.Priv != state.PrivM {
		cpu.Update.TrapMstatus &^= 1 << csr.MstatusMPRV
	}
}
