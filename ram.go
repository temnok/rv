package rv

func (cpu *CPU) accessRAM(addr, width int32) []byte {
	if addr -= ramBase; addr < 0 || addr+width > int32(len(cpu.ram)) {
		cpu.accessFault = true
		return nil
	}
	return cpu.ram[addr : addr+width]
}

func (cpu *CPU) readRAM(addr, width int32) int32 {
	val := int32(0)
	for i, b := range cpu.accessRAM(addr, width) {
		val |= int32(b) << (i << 3)
	}
	return val
}

func (cpu *CPU) writeRAM(addr, width int32, val int32) {
	ram := cpu.accessRAM(addr, width)
	for i := range ram {
		ram[i] = byte(val >> (i << 3))
	}
}
