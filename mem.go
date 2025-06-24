package rv

func (cpu *CPU) memFetch(virtAddr int32, data *int32) bool {
	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessExecute) {
		return false
	}

	if !cpu.bus.read(physAddr, data) {
		cpu.trap(ExceptionInstructionAccessFault)
		return false
	}

	if isWordAligned := virtAddr&3 == 0; isWordAligned {
		return true
	}

	*data = int32(uint32(*data) >> 16)
	if isCompressedInstruction := *data&3 != 3; isCompressedInstruction {
		return true
	}

	var hi int32
	if !cpu.bus.read(physAddr+1, &hi) {
		cpu.trap(ExceptionInstructionAccessFault)
		return false
	}

	*data |= hi << 16
	return true
}

func (cpu *CPU) memRead(virtAddr, width int32, data *int32) bool {
	if virtAddr&(width-1) != 0 {
		cpu.trap(ExceptionLoadAddressMisaligned)
		return false
	}

	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessRead) {
		return false
	}

	var word int32
	if !cpu.bus.read(physAddr, &word) {
		cpu.trap(ExceptionLoadAccessFault)
		return false
	}

	*data = word >> ((virtAddr & 3) << 3)
	return true
}

func (cpu *CPU) memWrite(virtAddr, width, data int32) bool {
	if virtAddr&(width-1) != 0 {
		cpu.trap(ExceptionStoreAMOAddressMisaligned)
		return false
	}

	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessWrite) {
		return false
	}

	if width < 4 {
		shift := (virtAddr & 3) << 3

		var old int32
		if !cpu.bus.read(physAddr, &old) {
			cpu.trap(ExceptionStoreAMOAccessFault)
			return false
		}

		if width == 1 {
			data = old&^(0xFF<<shift) | (data&0xFF)<<shift
		} else {
			data = old&^(0xFFFF<<shift) | (data&0xFFFF)<<shift
		}
	}

	if !cpu.bus.write(physAddr, data) {
		cpu.trap(ExceptionStoreAMOAccessFault)
		return false
	}

	return true
}
