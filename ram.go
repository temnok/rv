package rv

import "fmt"

type RAM struct {
	baseAddr Xint
	words    []Xint
}

func (ram *RAM) Init(baseAddr Xint, size int) {
	*ram = RAM{
		baseAddr: baseAddr,
		words:    make([]Xint, size/Xbytes),
	}
}

func (ram *RAM) Load(addr Xint, program []byte) {
	addr = (addr - ram.baseAddr) / Xbytes
	words := ram.words[addr : addr+Xint(len(program))/Xbytes+1]

	clear(words)
	for i, b := range program {
		shift := (i & (Xbytes - 1)) * 8
		words[i/Xbytes] |= Xint(b) << shift
	}
}

func (ram *RAM) access(addr Xint, data *Xint, width Xint, write bool) bool {
	i := (addr - ram.baseAddr) / Xbytes
	if i < 0 || i >= Xint(len(ram.words)) {
		return false
	}

	if width == Xbytes {
		if write {
			ram.words[i] = *data
		} else {
			*data = ram.words[i]
		}

		return true
	}

	if shift := (addr & (Xbytes - 1)) * 8; write {
		mask := Xint(-1) << (width * 8)
		ram.words[i] = (ram.words[i] &^ (^mask << shift)) | *data<<shift
	} else {
		*data = ram.words[i] >> shift
	}

	return true
}

func (ram *RAM) Dump(startAddr, byteCount Xint) {
	i0 := (startAddr - ram.baseAddr) / 8
	i1 := i0 + (byteCount+7)/8
	for i := i0; i < i1; i += 4 {
		fmt.Printf(
			"%016x: %016x %016x %016x %016x\r\n",
			ram.baseAddr+i*8,
			Xuint(ram.words[i]), Xuint(ram.words[i+1]),
			Xuint(ram.words[i+2]), Xuint(ram.words[i+3]),
		)
	}
}
