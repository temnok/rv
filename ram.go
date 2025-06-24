package rv

type RAM struct {
	baseAddr int32
	words    []int32
}

func (ram *RAM) Init(baseAddr int32, size int) {
	*ram = RAM{
		baseAddr: int32(uint32(baseAddr) >> 2),
		words:    make([]int32, (size+3)/4),
	}
}

func (ram *RAM) Load(addr int32, program []byte) {
	addr = int32(uint32(addr)/4) - ram.baseAddr
	words := ram.words[addr : addr+int32(len(program)+3)/4]

	clear(words)
	for i, b := range program {
		words[i/4] |= int32(b) << ((i & 3) * 8)
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
