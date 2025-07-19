package rv

type CPUUpdatedState struct {
	TrapEnter, TrapExit   bool
	TrapPriv, TrapPC      int
	TrapMstatus, TrapXepc int
	TrapXcause, TrapXtval int

	PC         int
	XReg, XVal int

	FReg, FVal int
	Fflags     int

	CReg, CVal int
	CRegPtr    *int

	Reserved     bool
	ReservedAddr int

	ICache Cache
}

func (cpu *CPU) updateState() {
	cpu.updateTimers()
	cpu.clearPendingInterrupts()

	up := &cpu.Updated

	if up.TrapEnter || up.TrapExit {
		cpu.Priv = up.TrapPriv
		cpu.PC = up.TrapPC
		cpu.CSR.Mstatus = up.TrapMstatus

		if up.TrapEnter {
			if up.TrapPriv == PrivM {
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

	if up.FReg >= 0 || up.CReg == Fflags || up.CReg == Frm || up.CReg == Fcsr {
		setBits(&cpu.CSR.Mstatus, MstatusFS, 2, FSdirty)
		cpu.CSR.Mstatus |= 1 << (cpu.XLen - 1) // set SD bit
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
	cpu.ICache = cpu.Updated.ICache
}

func (cpu *CPU) updateTimers() {
	if cpu.CSR.Cycle = cpu.xint(cpu.CSR.Cycle + 1); cpu.CSR.Cycle == 0 {
		cpu.CSR.Cycleh++
	}

	if cpu.CSR.Cycle%10_000 == 0 {
		if cpu.CSR.Time = cpu.xint(cpu.CSR.Time + 1); cpu.CSR.Time == 0 {
			cpu.CSR.Timeh++
		}
	}
}

func (cpu *CPU) clearPendingInterrupts() {
	cpu.CSR.Mip &^= 1<<MipSEI | 1<<MipMTI | 1<<MipMSI
}
