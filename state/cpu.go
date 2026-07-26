package state

type CPU struct {
	Priv        int
	PC          int
	X, F        [32]int
	CSR         CSR
	Reservation int

	TLB TLB
	RAM []int

	UARTInput  func() (byte, bool)
	UARTOutput func(byte) bool

	Update Update
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
