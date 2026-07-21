package mem

import (
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/translate"
	"github.com/temnok/rv/trap"
)

func Fetch(cpu *state.CPU, addr int) int {
	const xbytes = 8

	shift := addr & (xbytes - 1)
	virtAddr := addr &^ (xbytes - 1)

	var physAddr, val int
	if physAddr = translate.Sv(cpu, virtAddr, state.AccessFetch); trap.IsEntered(cpu) {
		return 0
	}

	val = cpu.RAM(physAddr, xbytes, false, 0)
	lo := val >> (shift * 8)
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= xbytes; fullyLoaded {
		return lo
	}

	virtAddr += xbytes
	physAddr += xbytes

	if pageMask := state.PageSize - 1; virtAddr&pageMask == 0 {
		if physAddr = translate.Sv(cpu, virtAddr, state.AccessFetch); trap.IsEntered(cpu) {
			return 0
		}
	}

	val = cpu.RAM(physAddr, xbytes, false, 0)

	return (val&0xffff)<<16 | lo&0xffff
}
