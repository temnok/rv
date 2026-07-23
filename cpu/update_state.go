package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func UpdateState(cpu *state.CPU) {
	up := &cpu.Update

	if up.Targets&state.UpdateXreg != 0 && up.Xreg != 0 {
		cpu.X[up.Xreg] = up.Xval
	}

	if up.Targets&state.UpdateFreg != 0 {
		cpu.F[up.Freg] = up.Fval
	}

	if up.Targets&state.UpdateCreg != 0 {
		csr.Update(cpu, up.Creg, up.Cval)
	}

	if up.Targets&state.UpdateFcsr != 0 {
		cpu.CSR.Fcsr = up.Fcsr
	}

	if up.Targets&state.UpdateMstatus != 0 {
		cpu.CSR.Mstatus = up.Mstatus
	}

	if up.Targets&state.UpdateEpc != 0 {
		if up.Priv == state.PrivM {
			cpu.CSR.Mepc = cpu.PC
		} else {
			cpu.CSR.Sepc = cpu.PC
		}
	}

	if up.Targets&state.UpdateCause != 0 {
		if up.Priv == state.PrivM {
			cpu.CSR.Mcause = up.Cause
		} else {
			cpu.CSR.Scause = up.Cause
		}
	}

	if up.Targets&state.UpdateTval != 0 {
		if up.Priv == state.PrivM {
			cpu.CSR.Mtval = up.Tval
		} else {
			cpu.CSR.Stval = up.Tval
		}
	}

	if up.Targets&state.UpdatePriv != 0 {
		cpu.Priv = up.Priv
	}

	if up.Targets&state.UpdateReservation != 0 {
		cpu.Reservation = up.Reservation
	}

	cpu.PC = up.PC

	updateCounters(cpu)
}

func updateCounters(cpu *state.CPU) {
	cpu.CSR.Mcycle++

	if uint(cpu.CSR.Mtime()) >= uint(cpu.CSR.Stimecmp) {
		cpu.CSR.Mip |= 1 << csr.MipSTI
	} else {
		cpu.CSR.Mip &^= 1 << csr.MipSTI
	}
}
