package ram

import "github.com/temnok/rv/state"

func ReadDword(cpu *state.CPU, addr int) int {
	if i := (addr - BaseAddr) >> 3; i >= 0 && i < len(cpu.RAM) {
		return cpu.RAM[i]
	}

	return 0
}

func Read(cpu *state.CPU, addr int, width int) int {
	if i := (addr - BaseAddr) >> 3; i >= 0 && i < len(cpu.RAM) {
		shift := (addr & 7) << 3
		mask := -1 << (width << 3)

		return cpu.RAM[i] >> shift &^ mask
	}

	return 0
}
