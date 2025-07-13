package rv

type CPUUpdatedState struct {
	TrapEnter, TrapExit   bool
	TrapPriv, TrapPC      int
	TrapMstatus, TrapXepc int
	TrapXcause, TrapXtval int

	PC                 int
	RegIndex, RegValue int

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

	if up.RegIndex != 0 {
		cpu.Reg[up.RegIndex] = up.RegValue
		up.RegIndex = 0
	}

	if up.CSRIndex != 0 {
		*up.CSRAddr = up.CSRValue
		up.CSRIndex = 0
	}

	cpu.Reserved = up.Reserved
	cpu.ReservedAddress = up.ReservedAddress
	cpu.ICache = cpu.Updated.ICache
}
