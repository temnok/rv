package state

type CPU struct {
	StaticState
	Update UpdatedState

	TLB TLB
	Bus Bus
}

func (cpu *CPU) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = val
}
