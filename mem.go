package rv

func (cpu *CPU) memFetch(virtAddr int, data *int) {
	const (
		xbytes   = 8
		pageMask = PageSize - 1
	)

	shift := virtAddr & (xbytes - 1)
	virtAddr &^= xbytes - 1

	var physAddr, val int
	if cpu.ICache.Hit(virtAddr) {
		physAddr, val = cpu.ICache.PhysAddr, cpu.ICache.Value
	} else {
		if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.isTrapped() {
			return
		}

		if !cpu.Bus.Read(physAddr, &val, xbytes) {
			cpu.trapEnter(ExceptionInstructionAccessFault, virtAddr)
			return
		}
	}

	lo := val >> (shift * 8)
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= xbytes; fullyLoaded {
		*data = lo

		cpu.Updated.ICache = Cache{
			VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

		return
	}

	virtAddr += xbytes
	physAddr += xbytes

	if virtAddr&pageMask == 0 {
		if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.isTrapped() {
			return
		}
	}

	if !cpu.Bus.Read(physAddr, &val, xbytes) {
		cpu.trapEnter(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	cpu.Updated.ICache = Cache{
		VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

	*data = val<<16 | bits(lo, 0, 16)
}

func (cpu *CPU) memRead(virtAddr int, data *int, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessRead); cpu.isTrapped() {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapEnter(ExceptionLoadAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Read(physAddr, data, width) {
		cpu.trapEnter(ExceptionLoadAccessFault, virtAddr)
	}
}

func (cpu *CPU) memWrite(virtAddr, data, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessWrite); cpu.isTrapped() {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.trapEnter(ExceptionStoreAMOAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Write(physAddr, data, width) {
		cpu.trapEnter(ExceptionStoreAMOAccessFault, virtAddr)
	}
}
