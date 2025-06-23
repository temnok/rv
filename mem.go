package rv

const ramBaseAddr = -1 << 31

func (cpu *CPU) memFetch(addr int32) int32 {
	if addr = cpu.translateSv32(addr, AccessExecute); cpu.trapped {
		return 0
	}

	mem := cpu.physMemAccess(addr, 2)
	if mem == nil {
		cpu.trap(ExceptionInstructionAccessFault)
		return 0
	}

	opcode := int32(mem[1])<<8 | int32(mem[0])
	if opcode&0b_11 == 0b_11 {
		mem = cpu.physMemAccess(addr+2, 2)
		if mem == nil {
			cpu.trap(ExceptionInstructionAccessFault)
			return 0
		}

		opcode |= int32(mem[1])<<24 | int32(mem[0])<<16
	}

	return opcode
}

func (cpu *CPU) memRead(addr, width int32) int32 {
	if addr = cpu.translateSv32(addr, AccessRead); cpu.trapped {
		return 0
	}

	return cpu.physMemRead(addr, width)
}

func (cpu *CPU) memWrite(addr, width, val int32) {
	if addr = cpu.translateSv32(addr, AccessWrite); cpu.trapped {
		return
	}

	mem := cpu.physMemAccess(addr, width)
	if mem == nil {
		cpu.trap(ExceptionStoreAMOAccessFault)
	}

	for i := range mem {
		mem[i] = byte(val >> (i << 3))
	}
}

func (cpu *CPU) physMemRead(addr, width int32) int32 {
	mem := cpu.physMemAccess(addr, width)
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

func (cpu *CPU) physMemAccess(addr, width int32) []byte {
	if addr -= ramBaseAddr; addr < 0 || addr+width > int32(len(cpu.mem)) {
		return nil
	}

	return cpu.mem[addr : addr+width]
}
