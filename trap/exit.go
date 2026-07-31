package trap

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#privstack
func Exit(cpu *state.CPU, retPriv int) {
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:mstatus_tsr_warl
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#otherpriv
	if retPriv > cpu.CSR.Priv || cpu.CSR.Priv == csr.PrivS && cpu.CSR.Mstatus>>csr.MstatusTSR&1 == 1 {
		Enter(cpu, IllegalIstruction, 0)
		return
	}

	cpu.Update.Targets = state.UpdatePC | state.UpdatePriv | state.UpdateMstatus

	status := cpu.CSR.Mstatus
	var nextStatus int

	if retPriv == csr.PrivM {
		cpu.Update.PC = cpu.CSR.Mepc
		cpu.Update.Priv = status >> csr.MstatusMPP & 3
		nextStatus = status&^(3<<csr.MstatusMPP) |
			1<<csr.MstatusMPIE | (status>>csr.MstatusMPIE&1)<<csr.MstatusMIE
	} else {
		cpu.Update.PC = cpu.CSR.Sepc
		cpu.Update.Priv = status >> csr.MstatusSPP & 1
		nextStatus = status&^(1<<csr.MstatusSPP) |
			1<<csr.MstatusSPIE | (status>>csr.MstatusSPIE&1)<<csr.MstatusSIE
	}

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-6-4-memory-privilege-in-mstatus-register
	if cpu.Update.Priv == csr.PrivM {
		cpu.Update.Mstatus = nextStatus
	} else {
		cpu.Update.Mstatus = nextStatus &^ (1 << csr.MstatusMPRV)
	}
}
