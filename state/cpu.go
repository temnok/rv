package state

type CPU struct {
	StaticState
	Update Update

	TLB TLB

	RAM     []int
	Devices Device
}

type Device func(addr int, width int, write bool, writeData int) int

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
