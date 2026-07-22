package state

import "github.com/temnok/rv/ram"

type CPU struct {
	StaticState
	Update UpdatedState

	TLB TLB

	RAM     Device
	Devices []Device

	InstrCount, CInstrCount, TrapCount int
}

type Device func(addr int, width int, write bool, writeData int) int

func (cpu *CPU) DeviceAtAddress(address int) Device {
	if address >= ram.BaseAddr {
		return cpu.RAM
	}

	if address < 0x1000_0000 {
		i := address >> 24
		if i < len(cpu.Devices) {
			return cpu.Devices[i]
		}
	}

	return nil
}

func (cpu *CPU) Xset(rd, val int) {
	cpu.Update.Targets |= UpdateXreg
	cpu.Update.Xreg = rd
	cpu.Update.Xval = val
}
