package state

type CPU struct {
	Fixed
	Static
	Update Updated

	TLB TLB
	Bus Bus
}

func (cpu *CPU) Xlen64() bool {
	return cpu.Xlen == 64
}

func (cpu *CPU) Xmask() int {
	return cpu.Xlen - 1
}

func (cpu *CPU) Xint(val int) int {
	if cpu.Xlen64() {
		return val
	}

	return int(int32(val))
}

func (cpu *CPU) Xuint(val int) uint {
	if cpu.Xlen64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *CPU) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = cpu.Xint(val)
}

func (cpu *CPU) XsetBool(rd int, val bool) {
	if val {
		cpu.Xset(rd, 1)
	} else {
		cpu.Xset(rd, 0)
	}
}

func (cpu *CPU) PCAddIf(c bool, imm int) {
	if c {
		cpu.Update.PC = cpu.Xint(cpu.PC + imm)
	}
}
