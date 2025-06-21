package rv

const ramBaseAddr = -1 << 31

func (cpu *CPU) memAccess(addr, width int32) []byte {
	if addr -= ramBaseAddr; addr < 0 || addr+width > int32(len(cpu.mem)) {
		cpu.memAccessFault = true
		return nil
	}

	return cpu.mem[addr : addr+width]
}

func (cpu *CPU) memRead(addr, width int32) int32 {
	val := int32(0)

	for i, b := range cpu.memAccess(addr, width) {
		val |= int32(b) << (i << 3)
	}

	return val
}

func (cpu *CPU) memWrite(addr, width, val int32) {
	bytes := cpu.memAccess(addr, width)

	for i := range bytes {
		bytes[i] = byte(val >> (i << 3))
	}
}
