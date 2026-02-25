package state

type CPU struct {
	StaticState
	Update UpdatedState

	TLB TLB
	Bus Bus

	InstrCount, CInstrCount, TrapCount int
}

func (cpu *CPU) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = val
}
