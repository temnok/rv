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

	var lo int32
	if !cpu.bus.read(physAddr, &lo) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return false
	}

	lo >>= (virtAddr & 3) * 8

	hiSize := virtAddr&3 + width - 4
	if hiSize <= 0 {
		*data = lo
		return true
	}

	loSize := width - hiSize
	virtAddr += loSize
	if !cpu.translateSv32(virtAddr, &physAddr, AccessRead) {
		return false
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return false
	}

	*data = hi<<(loSize*8) | lo&^(-1<<(loSize*8))
	return true
}

func (cpu *CPU) memWrite(virtAddr, width, data int32) bool {
	var physAddr int32
	if !cpu.translateSv32(virtAddr, &physAddr, AccessWrite) {
		return false
	}

	shift := virtAddr & 3
	if alignedWordWrite := width == 4 && shift == 0; alignedWordWrite {
		if !cpu.bus.write(physAddr, data) {
			cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
			return false
		}

		return true
	}

	var lo int32
	if !cpu.bus.read(physAddr, &lo) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return false
	}

	mask := (int32(0x100) << ((width - 1) * 8)) - 1
	lo &^= mask << (shift * 8)
	lo |= data << (shift * 8)

	if !cpu.bus.write(physAddr, lo) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return false
	}

	hiSize := shift + width - 4
	if hiSize <= 0 {
		return true
	}

	loSize := width - hiSize
	virtAddr += loSize
	if !cpu.translateSv32(virtAddr, &physAddr, AccessWrite) {
		return false
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return false
	}

	hi &= -1 << (hiSize * 8)
	hi |= int32(uint32(data) >> uint32(loSize*8))

	if !cpu.bus.write(physAddr, hi) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return false
	}

	return true
}
