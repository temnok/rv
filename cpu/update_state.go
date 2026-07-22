package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func UpdateState(cpu *state.CPU) {
	up := &cpu.Update

	if up.Targets&state.UpdatePriv != 0 {
		cpu.Priv = up.Priv
	}

	if up.Targets&state.UpdateMstatus != 0 {
		cpu.CSR.Mstatus = up.Mstatus
	}

	if up.Targets&state.UpdateReservation != 0 {
		cpu.Reservation = up.Reservation
	}

	if up.Targets&state.UpdateEpc != 0 {
		if up.Priv == state.PrivM {
			cpu.CSR.Mepc = cpu.PC
		} else {
			cpu.CSR.Sepc = cpu.PC
		}
	}

	if up.Targets&state.UpdateXreg != 0 && up.Xreg != 0 {
		cpu.X[up.Xreg] = up.Xval
	}

	if up.Targets&state.UpdateFreg != 0 {
		cpu.F[up.Freg] = up.Fval
	}

	if up.Targets&state.UpdateCreg != 0 {
		csr.Update(cpu, up.Creg, up.Cval)
	}

	if up.Targets&state.UpdateFflags != 0 {
		cpu.CSR.Fcsr |= up.Fflags
	}

	cpu.PC = up.PC

	up.Targets = 0

	if up.TrapEnter {
		if up.Priv == state.PrivM {
			cpu.CSR.Mcause = up.TrapXcause
			cpu.CSR.Mtval = up.TrapXtval
		} else {
			cpu.CSR.Scause = up.TrapXcause
			cpu.CSR.Stval = up.TrapXtval
		}

		up.TrapEnter = false
		return
	}

	creg := up.Creg

	sd := csr.MstatusSD64
	if up.Targets&state.UpdateFreg != 0 ||
		up.Targets&state.UpdateCreg != 0 && (creg == csr.Fflags || creg == csr.Frm || creg == csr.Fcsr) {
		cpu.CSR.Mstatus &^= 0b_11 << csr.MstatusFS
		cpu.CSR.Mstatus |= csr.FSdirty << csr.MstatusFS
		cpu.CSR.Mstatus |= 1 << sd
	}

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
