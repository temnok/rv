package mem

import (
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Fetch(cpu *state.CPU, addr int) int {
	const (
		xbytes   = 8
		pageSize = 0x1000
	)

	shift := addr & (xbytes - 1)
	virtAddr := addr &^ (xbytes - 1)

	var physAddr, val int
	if physAddr = translateSv39(cpu, virtAddr, accessFetch); trap.IsEntered(cpu) {
		return -1
	}

	val = ram.Read8(cpu, physAddr)
	lo := val >> (shift * 8)
	isCompressedInstruction := lo&3 != 3

	if fullyLoaded := isCompressedInstruction || shift+4 <= xbytes; fullyLoaded {
		return lo
	}

	virtAddr += xbytes
	physAddr += xbytes

	if pageMask := pageSize - 1; virtAddr&pageMask == 0 {
		if physAddr = translateSv39(cpu, virtAddr, accessFetch); trap.IsEntered(cpu) {
			return -1
		}
	}

	val = ram.Read8(cpu, physAddr)

	return (val&0xffff)<<16 | lo&0xffff
}
