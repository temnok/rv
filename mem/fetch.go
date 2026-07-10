package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Fetch(cpu *state.CPU, addr int, data *int) {
	const xbytes = 8

	shift := addr & (xbytes - 1)
	virtAddr := addr &^ (xbytes - 1)

	var physAddr, val int
	var ok bool
	if cpu.ICache.Hit(virtAddr) {
		physAddr, val = cpu.ICache.PhysAddr, cpu.ICache.Value
	} else {
		if translate.Sv(cpu, virtAddr, &physAddr, state.AccessExecute); trap.IsEntered(cpu) {
			return
		}

		if val, ok = cpu.Bus.Read(physAddr, xbytes); !ok {
			trap.Enter(cpu, trap.InstructionAccessFault, addr)
			return
		}
	}

	lo := val >> (shift * 8)
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= xbytes; fullyLoaded {
		*data = lo

		cpu.Update.ICache = state.Cache{VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

		return
	}

	virtAddr += xbytes
	physAddr += xbytes

	if pageMask := state.PageSize - 1; virtAddr&pageMask == 0 {
		if translate.Sv(cpu, virtAddr, &physAddr, state.AccessExecute); trap.IsEntered(cpu) {
			return
		}
	}

	if val, ok = cpu.Bus.Read(physAddr, xbytes); !ok {
		trap.Enter(cpu, trap.InstructionAccessFault, virtAddr)
		return
	}

	cpu.Update.ICache = state.Cache{VirtAddr: virtAddr, PhysAddr: physAddr, Value: val}

	*data = (val&0xffff)<<16 | lo&0xffff
}
