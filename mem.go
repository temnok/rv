package rv

const ramBaseAddr = -1 << 31

func (cpu *CPU) memAccess(addr, width int32) []byte {
	if addr -= ramBaseAddr; addr < 0 || addr+width > int32(len(cpu.mem)) {
		return nil
	}

	return cpu.mem[addr : addr+width]
}

func (cpu *CPU) memFetch(addr int32) int32 {
	mem := cpu.memAccess(addr, 4)
	if mem == nil {
		cpu.trap(ExceptionInstructionAccessFault)
		return 0
	}

	val := int32(0)

	for i, b := range mem {
		val |= int32(b) << (i << 3)
	}

	return val
}

func (cpu *CPU) memRead(addr, width int32) int32 {
	mem := cpu.memAccess(addr, width)
	if mem == nil {
		cpu.trap(ExceptionLoadAccessFault)
		return 0
	}

	val := int32(0)

	for i, b := range mem {
		val |= int32(b) << (i << 3)
	}

	return val
}

func (cpu *CPU) memWrite(addr, width, val int32) {
	mem := cpu.memAccess(addr, width)
	if mem == nil {
		cpu.trap(ExceptionStoreAMOAccessFault)
	}

	for i := range mem {
		mem[i] = byte(val >> (i << 3))
	}
}
