package state

import "github.com/temnok/rv/arch"

type CPU struct {
	StaticState
	Update UpdatedState

	TLB TLB
	Bus Bus
}

func (cpu *CPU) LenIs64() bool {
	return true
}

func (cpu *CPU) Mask() int {
	return arch.XMask
}

func (cpu *CPU) Int(val int) int {
	if cpu.LenIs64() {
		return val
	}

	return int(int32(val))
}

func (cpu *CPU) Uint(val int) uint {
	if cpu.LenIs64() {
		return uint(val)
	}

	return uint(uint32(val))
}

func (cpu *CPU) Xset(rd, val int) {
	cpu.Update.XReg = rd
	cpu.Update.XVal = cpu.Int(val)
}
