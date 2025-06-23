package rv

func (cpu *CPU) memFetch(addr int32, data *int32) bool {
	if !cpu.translateSv32(&addr, AccessExecute) {
		return false
	}

	if !cpu.bus.read(addr&^3, 4, data) {
		cpu.trap(ExceptionInstructionAccessFault)
		return false
	}

	if isWordAligned := addr&3 == 0; isWordAligned {
		return true
	}

	*data = int32(uint32(*data) >> 16)
	if isCompressedInstruction := *data&3 != 3; isCompressedInstruction {
		return true
	}

	var hi int32
	if !cpu.bus.read(addr+2, 4, &hi) {
		cpu.trap(ExceptionInstructionAccessFault)
		return false
	}

	*data |= hi << 16
	return true
}

func (cpu *CPU) memRead(addr, width int32, data *int32) bool {
	if addr&(width-1) != 0 {
		cpu.trap(ExceptionLoadAddressMisaligned)
		return false
	}

	if !cpu.translateSv32(&addr, AccessRead) {
		return false
	}

	var word int32
	if !cpu.bus.read(addr&^3, 4, &word) {
		cpu.trap(ExceptionLoadAccessFault)
		return false
	}

	*data = word >> ((addr & 3) << 3)
	return true
}

func (cpu *CPU) memWrite(addr, width, data int32) bool {
	if addr&(width-1) != 0 {
		cpu.trap(ExceptionStoreAMOAddressMisaligned)
		return false
	}

	if !cpu.translateSv32(&addr, AccessWrite) {
		return false
	}

	if width < 4 {
		shift := (addr & 3) << 3
		addr &= ^3

		var old int32
		if !cpu.bus.read(addr, 4, &old) {
			cpu.trap(ExceptionStoreAMOAccessFault)
			return false
		}

		if width == 1 {
			data = old&^(0xFF<<shift) | (data&0xFF)<<shift
		} else {
			data = old&^(0xFFFF<<shift) | (data&0xFFFF)<<shift
		}
	}

	if !cpu.bus.write(addr, 4, data) {
		cpu.trap(ExceptionStoreAMOAccessFault)
		return false
	}

	return true
}
