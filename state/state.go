package state

type State struct {
	Permanent
	Static
	Update Updated
}

func (cpu *State) XLen64() bool {
	return cpu.XLen == 64
}

func (cpu *State) Xint(val int) int {
	if cpu.XLen64() {
		return val
	}

	return int(int32(val))
}

func (cpu *State) Xuint(val int) uint {
	if cpu.XLen64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *State) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = cpu.Xint(val)
}

func (cpu *State) PCAdd(imm int) {
	cpu.Update.PC = cpu.Xint(cpu.PC + imm)
}
