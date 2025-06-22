package rv

const ramBaseAddr = -1 << 31

func (cpu *CPU) memAccess(addr, width int32) []byte {
	if addr -= ramBaseAddr; addr < 0 || addr+width > int32(len(cpu.mem)) {
		return nil
	}

	return cpu.mem[addr : addr+width]
}

func (cpu *CPU) memRead(addr, width int32, out *int32) bool {
	mem := cpu.memAccess(addr, width)
	if mem == nil {
		return false
	}

	val := int32(0)

	for i, b := range mem {
		val |= int32(b) << (i << 3)
	}

	*out = val
	return true
}

func (cpu *CPU) memWrite(addr, width, val int32) bool {
	mem := cpu.memAccess(addr, width)
	if mem == nil {
		return false
	}

	for i := range mem {
		mem[i] = byte(val >> (i << 3))
	}

	return true
}
