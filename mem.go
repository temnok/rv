package rv

func (cpu *CPU) memFetch(virtAddr int, data *int) {
	shift := virtAddr & (Xbytes - 1)
	virtAddr &^= Xbytes - 1

	var physAddr, lo int
	if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.isTrapped {
		return
	}

	if !cpu.bus.read(physAddr, &lo, Xbytes) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	lo >>= shift * 8
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= Xbytes; fullyLoaded {
		*data = lo
		return
	}

	virtAddr += Xbytes
	if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.isTrapped {
		return
	}

	var hi int
	if !cpu.bus.read(physAddr, &hi, Xbytes) {
		cpu.trapWithTval(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	*data = hi<<16 | bits(lo, 0, 16)
}

func (cpu *CPU) memRead(virtAddr int, data *int, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessRead); cpu.isTrapped {
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

func (cpu *CPU) memWrite(virtAddr, data, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessWrite); cpu.isTrapped {
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
