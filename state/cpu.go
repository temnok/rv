package state

type CPU struct {
	StaticState
	Update Update

	TLB TLB

	RAM []int

	UARTInput  func() (byte, bool)
	UARTOutput func(byte) bool
}

func (cpu *CPU) Xset(reg, val int) {
	cpu.Update.Targets |= UpdateXreg
	cpu.Update.Xreg = reg
	cpu.Update.Xval = val
}

func (cpu *CPU) Fset(reg, val int) {
	cpu.Update.Targets |= UpdateFreg | UpdateMstatus
	cpu.Update.Freg = reg
	cpu.Update.Fval = val
	cpu.Update.Mstatus = cpu.CSR.Mstatus | MstatusDirtyMask
}

func (cpu *CPU) Cset(reg, val int) {
	cpu.Update.Targets |= UpdateCreg
	cpu.Update.Creg = reg
	cpu.Update.Cval = val

	if reg >= 1 && reg <= 3 { // Fflags, Frm, Fcsr
		cpu.Update.Targets |= UpdateMstatus
		cpu.Update.Mstatus = cpu.CSR.Mstatus | MstatusDirtyMask
	}
}
