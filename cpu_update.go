package rv

type CPUUpdatedState struct {
	TrapEnter, TrapExit   bool
	TrapPriv, TrapPC      int
	TrapMstatus, TrapXepc int
	TrapXcause, TrapXtval int

	PC                 int
	RegIndex, RegValue int

	FRegUpdated          bool
	FRegIndex, FRegValue int
	Fflags               int

	CSRAddr            *int
	CSRIndex, CSRValue int

	Reserved        bool
	ReservedAddress int

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

	if up.RegIndex > 0 {
		cpu.Reg[up.RegIndex] = up.RegValue
		up.RegIndex = -1
	}

	if up.FRegIndex >= 0 || up.CSRIndex == Fflags || up.CSRIndex == Frm || up.CSRIndex == Fcsr {
		setBits(&cpu.CSR.Mstatus, MstatusFS, 2, FSdirty)
		cpu.CSR.Mstatus |= 1 << (cpu.XLen - 1) // set SD bit
	}

	if up.FRegIndex >= 0 {
		cpu.FReg[up.FRegIndex] = up.FRegValue
		up.FRegIndex = -1
	}

	if up.Fflags != 0 {
		cpu.CSR.Fcsr |= up.Fflags
		up.Fflags = 0
	}

	if up.CSRIndex >= 0 {
		*up.CSRAddr = up.CSRValue
		up.CSRIndex = -1
	}

	cpu.Reserved = up.Reserved
	cpu.ReservedAddress = up.ReservedAddress
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
