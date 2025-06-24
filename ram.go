package rv

type RAM struct {
	baseAddr int32
	words    []int32
}

func (ram *RAM) Init(baseAddr int32, size int, program []byte) {
	*ram = RAM{
		baseAddr: int32(uint32(baseAddr) >> 2),
		words:    make([]int32, size/4),
	}

	for i, b := range program {
		ram.words[i>>2] |= int32(b) << ((i & 3) << 3)
	}
}

func (ram *RAM) access(addr int32, data *int32, write bool) bool {
	if addr -= ram.baseAddr; addr < 0 || addr >= int32(len(ram.words)) {
		return false
	}

	if write {
		ram.words[addr] = *data
	} else {
		*data = ram.words[addr]
	}

	return true
}
