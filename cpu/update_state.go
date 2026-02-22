package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func updateState(cpu *state.CPU) {
	updateTimers(cpu)
	clearPendingInterrupts(cpu)

	up := &cpu.Update

	if up.TrapEnter || up.TrapExit {
		cpu.Priv = up.TrapPriv
		cpu.PC = up.TrapPC
		cpu.CSR.Mstatus = up.TrapMstatus

		if up.TrapEnter {
			if up.TrapPriv == state.PrivM {
				cpu.CSR.Mepc = up.TrapXepc
				cpu.CSR.Mcause = up.TrapXcause
				cpu.CSR.Mtval = up.TrapXtval
			} else {
				cpu.CSR.Sepc = up.TrapXepc
				cpu.CSR.Scause = up.TrapXcause
				cpu.CSR.Stval = up.TrapXtval
			}
		}

		up.TrapEnter = false
		up.TrapExit = false
		return
	}

	cpu.PC = up.PC

	if up.XReg > 0 {
		cpu.X[up.XReg] = up.XVal
		up.XReg = -1
	}

	sd := csr.MstatusSD64
	if up.FReg >= 0 || up.CReg == csr.Fflags || up.CReg == csr.Frm || up.CReg == csr.Fcsr {
		cpu.CSR.Mstatus &^= 0b_11 << csr.MstatusFS
		cpu.CSR.Mstatus |= csr.FSdirty << csr.MstatusFS
		cpu.CSR.Mstatus |= 1 << sd
	}

	if up.FReg >= 0 {
		cpu.F[up.FReg] = up.FVal
		up.FReg = -1
	}

	if up.Fflags != 0 {
		cpu.CSR.Fcsr |= up.Fflags
		up.Fflags = 0
	}

	if up.CReg >= 0 {
		*up.CRegPtr = up.CVal
		up.CReg = -1
	}

	cpu.Reserved = up.Reserved
	cpu.ReservedAddr = up.ReservedAddr
	cpu.ICache = cpu.Update.ICache
}

func updateTimers(cpu *state.CPU) {
	if cpu.CSR.Cycle = cpu.Int(cpu.CSR.Cycle + 1); cpu.CSR.Cycle == 0 {
		cpu.CSR.Cycleh++
	}

	if cpu.CSR.Cycle%10_000 == 0 {
		if cpu.CSR.Time = cpu.Int(cpu.CSR.Time + 1); cpu.CSR.Time == 0 {
			cpu.CSR.Timeh++
		}
	}
}

func clearPendingInterrupts(cpu *state.CPU) {
	cpu.CSR.Mip &^= 1<<csr.MipSEI | 1<<csr.MipMTI | 1<<csr.MipMSI
}
