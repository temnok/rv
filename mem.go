package rv

func (cpu *CPU) memFetch(virtAddr int32, data *int32) bool {
	var physAddr, lo int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessExecute) {
		return false
	}

	if !cpu.bus.read(physAddr, &lo) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return false
	}

	if isWordAligned := virtAddr&3 == 0; isWordAligned {
		*data = lo
		return true
	}

	lo = int32(uint32(lo) >> 16)

	if isCompressedInstruction := bits(lo, 0, 2) != 0b_11; isCompressedInstruction {
		*data = lo
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

	*data = hi<<16 | lo
	return true
}

func (cpu *CPU) memRead(virtAddr, width int32, data *int32) bool {
	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessRead) {
		return false
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionLoadAddressMisaligned, virtAddr)
		return false
	}

	var word int32
	if !cpu.bus.read(physAddr, &word) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return false
	}

	*data = word >> ((virtAddr & 3) * 8)
	return true
}

func (cpu *CPU) memWrite(virtAddr, width, data int32) bool {
	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessWrite) {
		return false
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionStoreAMOAddressMisaligned, virtAddr)
		return false
	}

	if width < 4 {
		shift := (virtAddr & 3) * 8

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
