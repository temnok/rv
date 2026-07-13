package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func UpdateState(cpu *state.CPU) {
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

	creg := up.CReg

	sd := csr.MstatusSD64
	if up.FReg >= 0 || creg == csr.Fflags || creg == csr.Frm || creg == csr.Fcsr {
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

	if creg >= 0 {
		*up.CRegPtr = up.CVal
		up.CReg = -1
	}

	cpu.Reserved = up.Reserved
	cpu.ReservedAddr = up.ReservedAddr
	cpu.ICache = cpu.Update.ICache

	updateCounters(cpu, creg)
}

func updateCounters(cpu *state.CPU, creg int) {
	//if cpu.CSR.Mcountinhibit&(1<<csr.McountinhibitCY) == 0 {
	cpu.CSR.Mcycle++

	//if cpu.CSR.Mcountinhibit&(1<<csr.McountinhibitIR) == 0 {
	cpu.CSR.Minstret++

	checkStimecmp := creg == csr.Stimecmp

	if cpu.CSR.Mcycle%20_000 == 0 {
		cpu.CSR.Mtime++

		for _, c := range cpu.CSR.TimerCallbacks {
			c()
		}

		checkStimecmp = true
	}

	if checkStimecmp {
		if uint(cpu.CSR.Mtime) >= uint(cpu.CSR.Stimecmp) {
			cpu.CSR.Mip |= 1 << csr.MipSTI
		} else {
			cpu.CSR.Mip &^= 1 << csr.MipSTI
		}
	}
}
