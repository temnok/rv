package rv

type RAM struct {
	words []int32
}

const (
	ramBaseAddr = -1 << 31
)

func (ram *RAM) init(size int, program []byte) {
	*ram = RAM{
		words: make([]int32, (size+3)/4),
	}

	for i, b := range program[:min(len(program), size)] {
		ram.words[i>>2] |= int32(b) << ((i & 3) << 3)
	}
}

func (ram *RAM) access(addr int32, data *int32, write bool) bool {
	ramBase := int32(ramBaseAddr)

	if addr -= int32(uint32(ramBase) >> 2); addr < 0 || addr >= int32(len(ram.words)) {
		return false
	}

	if write {
		ram.words[addr] = *data
	} else {
		*data = ram.words[addr]
	}

	return true
}
