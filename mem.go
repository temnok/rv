package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func (cpu *CPU) memFetch(addr int, data *int) {
	const (
		xbytes   = 8
		pageMask = PageSize - 1
	)

	shift := addr & (xbytes - 1)
	virtAddr := addr &^ (xbytes - 1)

	var physAddr, val int
	if cpu.ICache.Hit(virtAddr) {
		physAddr, val = cpu.ICache.PhysAddr, cpu.ICache.Value
	} else {
		if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.IsTrapped() {
			return
		}

		if !cpu.Bus.Read(physAddr, &val, xbytes) {
			cpu.TrapEnter(ExceptionInstructionAccessFault, addr)
			return
		}
	}

	lo := val >> (shift * 8)
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= xbytes; fullyLoaded {
		*data = lo

		cpu.Update.ICache = state.Cache{
			VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

		return
	}

	virtAddr += xbytes
	physAddr += xbytes

	if virtAddr&pageMask == 0 {
		if cpu.translateSv(virtAddr, &physAddr, AccessExecute); cpu.IsTrapped() {
			return
		}
	}

	if !cpu.Bus.Read(physAddr, &val, xbytes) {
		cpu.TrapEnter(ExceptionInstructionAccessFault, virtAddr)
		return
	}

	cpu.Update.ICache = state.Cache{
		VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

	*data = val<<16 | bi.Ts(lo, 0, 16)
}

func (cpu *CPU) memRead(virtAddr int, data *int, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessRead); cpu.IsTrapped() {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.TrapEnter(ExceptionLoadAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Read(physAddr, data, width) {
		cpu.TrapEnter(ExceptionLoadAccessFault, virtAddr)
	}
}

func (cpu *CPU) memWrite(virtAddr, data, width int) {
	var physAddr int
	if cpu.translateSv(virtAddr, &physAddr, AccessWrite); cpu.IsTrapped() {
		return
	}

	if virtAddr&(width-1) != 0 {
		cpu.TrapEnter(ExceptionStoreAMOAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Write(physAddr, data, width) {
		cpu.TrapEnter(ExceptionStoreAMOAccessFault, virtAddr)
	}
}
