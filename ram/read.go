package ram

func Read8(ram []int, addr int) int {
	if i := (addr - BaseAddr) >> 3; i >= 0 && i < len(ram) {
		return ram[i]
	}

	return 0
}

func Read(ram []int, addr int, width int) int {
	if i := (addr - BaseAddr) >> 3; i >= 0 && i < len(ram) {
		shift := (addr & 7) << 3
		mask := -1 << (width << 3)

		return ram[i] >> shift &^ mask
	}

	return 0
}
