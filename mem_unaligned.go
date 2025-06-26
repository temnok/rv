package rv

func (cpu *CPU) memRead_unaligned(virtAddr int32, data *int32, width int32) {
	var physAddr int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessRead); cpu.isTrapped {
		return
	}

	var lo int32
	if !cpu.bus.read(physAddr, &lo) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	lo >>= (virtAddr & 3) * 8

	hiSize := virtAddr&3 + width - 4
	if hiSize <= 0 {
		*data = lo
		return
	}

	loSize := width - hiSize
	virtAddr += loSize
	if cpu.translateSv32(virtAddr, &physAddr, AccessRead); cpu.isTrapped {
		return
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi) {
		cpu.trapWithTval(ExceptionLoadAccessFault, virtAddr)
		return
	}

	*data = hi<<(loSize*8) | lo&^(-1<<(loSize*8))
}

func (cpu *CPU) memWrite_unaligned(virtAddr, data, width int32) {
	var physAddr int32
	if cpu.translateSv32(virtAddr, &physAddr, AccessWrite); cpu.isTrapped {
		return
	}

	shift := virtAddr & 3
	if alignedWordWrite := width == 4 && shift == 0; alignedWordWrite {
		if !cpu.bus.write(physAddr, data) {
			cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		}

		return
	}

	var lo int32
	if !cpu.bus.read(physAddr, &lo) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return
	}

	mask := (int32(0x100) << ((width - 1) * 8)) - 1
	lo &^= mask << (shift * 8)
	lo |= data << (shift * 8)

	if !cpu.bus.write(physAddr, lo) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return
	}

	hiSize := shift + width - 4
	if hiSize <= 0 {
		return
	}

	loSize := width - hiSize
	virtAddr += loSize
	if cpu.translateSv32(virtAddr, &physAddr, AccessWrite); cpu.isTrapped {
		return
	}

	var hi int32
	if !cpu.bus.read(physAddr, &hi) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
		return
	}

	hi &= -1 << (hiSize * 8)
	hi |= int32(uint32(data) >> uint32(loSize*8))

	if !cpu.bus.write(physAddr, hi) {
		cpu.trapWithTval(ExceptionStoreAMOAccessFault, virtAddr)
	}
}
