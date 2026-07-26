package ram

import "github.com/temnok/rv/state"

func Write(cpu *state.CPU, addr int, width int, val int) {
	i := (addr - BaseAddr) >> 3
	if i < 0 || i >= len(cpu.RAM) {
		return
	}

	if width == 8 {
		cpu.RAM[i] = val
		return
	}

	shift := (addr & 7) << 3
	mask := 1<<(width<<3) - 1

	cpu.Update.Targets |= state.UpdateRAM
	cpu.Update.RAMPos = i
	cpu.Update.RAMVal = cpu.RAM[i]&^(mask<<shift) | (val&mask)<<shift
}
