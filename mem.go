package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
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
		if translate.Sv(cpu.State, virtAddr, &physAddr, AccessExecute); trap.IsEntered(cpu.State) {
			return
		}

		if !cpu.Bus.Read(physAddr, &val, xbytes) {
			trap.Enter(cpu.State, trap.InstructionAccessFault, addr)
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
		if translate.Sv(cpu.State, virtAddr, &physAddr, AccessExecute); trap.IsEntered(cpu.State) {
			return
		}
	}

	if !cpu.Bus.Read(physAddr, &val, xbytes) {
		trap.Enter(cpu.State, trap.InstructionAccessFault, virtAddr)
		return
	}

	cpu.Update.ICache = state.Cache{
		VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

	*data = val<<16 | bi.Ts(lo, 0, 16)
}

func (cpu *CPU) memRead(virtAddr int, data *int, width int) {
	var physAddr int
	if translate.Sv(cpu.State, virtAddr, &physAddr, AccessRead); trap.IsEntered(cpu.State) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu.State, trap.LoadAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Read(physAddr, data, width) {
		trap.Enter(cpu.State, trap.LoadAccessFault, virtAddr)
	}
}

func (cpu *CPU) memWrite(virtAddr, data, width int) {
	var physAddr int
	if translate.Sv(cpu.State, virtAddr, &physAddr, AccessWrite); trap.IsEntered(cpu.State) {
		return
	}

	if virtAddr&(width-1) != 0 {
		trap.Enter(cpu.State, trap.StoreAMOAddressMisaligned, virtAddr)
		return
	}

	if !cpu.Bus.Write(physAddr, data, width) {
		trap.Enter(cpu.State, trap.StoreAMOAccessFault, virtAddr)
	}
}
