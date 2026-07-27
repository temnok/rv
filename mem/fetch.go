package mem

import (
	"github.com/temnok/rv/ram"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Fetch(cpu *state.CPU, addr int) int {
	va := addr &^ 7

	pa, shift := translateSv39(cpu, va, accessFetch)
	if trap.IsEntered(cpu) {
		return -1
	}

	dword := ram.ReadDword(cpu, pa)
	offset := addr & 7
	opcode := dword >> (offset << 3)
	isCompressedInstruction := opcode&3 != 3

	if fullyLoaded := isCompressedInstruction || offset+4 <= 8; fullyLoaded {
		return opcode
	}

	if va>>shift == (va+8)>>shift {
		pa += 8
	} else {
		pa, _ = translateSv39(cpu, va+8, accessFetch)
		if trap.IsEntered(cpu) {
			return -1
		}
	}

	dword = ram.ReadDword(cpu, pa)

	return (dword&0xffff)<<16 | opcode&0xffff
}
