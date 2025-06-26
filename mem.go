package rv

func (cpu *CPU) memFetch(virtAddr int32, data *int32) {
	var physAddr, lo int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessExecute); cpu.isTrapped {
		return
	}

	if !cpu.bus.read(physAddr, &lo) {
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
	if !cpu.bus.read(physAddr, &hi) {
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

	var word int32
	if !cpu.bus.read(physAddr, &word) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	shift := (virtAddr & 3) * 8
	*data = word >> shift
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

	if width < 4 {
		var old int32
		if !cpu.bus.read(physAddr, &old) {
			cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
			return
		}

		shift := (virtAddr & 3) * 8
		mask := ^(int32(-1) << (width * 8)) << shift
		data = old&^mask | data<<shift
	}

	if !cpu.bus.write(physAddr, data) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
	}
}
