package state

type State struct {
	Fixed
	Static
	Update Updated
}

func (cpu *State) Xlen64() bool {
	return cpu.Xlen == 64
}

func (cpu *State) Xmask() int {
	return cpu.Xlen - 1
}

func (cpu *State) Xint(val int) int {
	if cpu.Xlen64() {
		return val
	}

	return int(int32(val))
}

func (cpu *State) Xuint(val int) uint {
	if cpu.Xlen64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *State) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = cpu.Xint(val)
}

func (cpu *State) XsetBool(rd int, val bool) {
	if val {
		cpu.Xset(rd, 1)
	} else {
		cpu.Xset(rd, 0)
	}
}

func (cpu *State) PCAddIf(c bool, imm int) {
	if c {
		cpu.Update.PC = cpu.Xint(cpu.PC + imm)
	}
}
