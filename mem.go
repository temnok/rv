package rv

func (cpu *CPU) memFetch(virtAddr int32, data *int32) bool {
	var physAddr, word int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessExecute) {
		return false
	}

	if !cpu.bus.read(physAddr, &word) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return false
	}

	if isWordAligned := virtAddr&3 == 0; isWordAligned {
		*data = word
		return true
	}

	word = int32(uint32(word) >> 16)

	if isCompressedInstruction := word&3 != 3; isCompressedInstruction {
		*data = word
		return true
	}

	virtAddr += 2
	if !cpu.translateSv32(virtAddr, &physAddr, AccessExecute) {
		return false
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return false
	}

	*data = hi<<16 | word
	return true
}

func (cpu *CPU) memRead(virtAddr, width int32, data *int32) bool {
	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionLoadAddressMisaligned, virtAddr)
		return false
	}

	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessRead) {
		return false
	}

	var word int32
	if !cpu.bus.read(physAddr, &word) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return false
	}

	*data = word >> ((virtAddr & 3) << 3)
	return true
}

func (cpu *CPU) memWrite(virtAddr, width, data int32) bool {
	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionStoreAMOAddressMisaligned, virtAddr)
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
			cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
			return false
		}

		if width == 1 {
			data = old&^(0xFF<<shift) | (data&0xFF)<<shift
		} else {
			data = old&^(0xFFFF<<shift) | (data&0xFFFF)<<shift
		}
	}

	if !cpu.bus.write(physAddr, data) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return false
	}

	return true
}
