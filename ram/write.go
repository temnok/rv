package ram

func Write(ram []int, addr int, width int, val int) {
	i := (addr - BaseAddr) >> 3
	if i < 0 || i >= len(ram) {
		return
	}

	if width == 8 {
		ram[i] = val
		return
	}

	shift := (addr & 7) << 3
	mask := 1<<(width<<3) - 1

	ram[i] = ram[i]&^(mask<<shift) | (val&mask)<<shift
}
