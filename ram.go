package rv

type RAM struct {
	baseAddr int32
	words    []int32
}

func (ram *RAM) Init(baseAddr int32, size int) {
	*ram = RAM{
		baseAddr: baseAddr,
		words:    make([]int32, size/4),
	}
}

func (ram *RAM) Load(addr int32, program []byte) {
	addr = (addr - ram.baseAddr) / 4
	words := ram.words[addr : addr+int32(len(program))/4+1]

	clear(words)
	for i, b := range program {
		words[i/4] |= int32(b) << ((i & 3) * 8)
	}
}

func (ram *RAM) access(addr int32, data *int32, width int32, write bool) bool {
	i := (addr - ram.baseAddr) / 4
	if i < 0 || i >= int32(len(ram.words)) {
		return false
	}

	if width == 4 {
		if write {
			ram.words[i] = *data
		} else {
			*data = ram.words[i]
		}

		return true
	}

	if shift := (addr & 3) * 8; write {
		mask := int32(-1) << (width * 8)
		ram.words[i] = (ram.words[i] &^ (^mask << shift)) | *data<<shift
	} else {
		*data = ram.words[i] >> shift
	}

	return true
}
