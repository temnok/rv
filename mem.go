package rv

func (cpu *CPU) memFetch(virtAddr int32, data *int32) {
	var physAddr, lo int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessExecute); cpu.isTrapped {
		return
	}

	if !cpu.bus.read(physAddr, &lo, 4) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	if isWordAligned := virtAddr&3 == 0; isWordAligned {
		*data = lo
		return
	}

	lo = int32(uint32(lo) >> 16)

	if isCompressedInstruction := bits(lo, 0, 2) != 0b_11; isCompressedInstruction {
		*data = lo
		return
	}

	virtAddr += 2
	if cpu.translateSv32(virtAddr, &physAddr, AccessExecute); cpu.isTrapped {
		return
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi, 4) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	*data = hi<<16 | lo
}

func (cpu *CPU) memRead(virtAddr int32, data *int32, width int32) {
	var physAddr int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessRead); cpu.isTrapped {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionLoadAddressMisaligned, virtAddr)
		return
	}

	if !cpu.bus.read(physAddr, data, width) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
	}
}

func (cpu *CPU) memWrite(virtAddr, data, width int32) {
	var physAddr int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessWrite); cpu.isTrapped {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapWithTval(ExceptionStoreAMOAddressMisaligned, virtAddr)
		return
	}

	if !cpu.bus.write(physAddr, data, width) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
	}
}
