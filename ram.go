package rv

import "fmt"

type RAM struct {
	cpu      *CPU
	baseAddr int
	words    []int
}

func (ram *RAM) Init(cpu *CPU, baseAddr int, size int) {
	*ram = RAM{
		cpu:      cpu,
		baseAddr: cpu.xint(baseAddr),
		words:    make([]int, size/8),
	}
}

func (ram *RAM) Load(addr int, program []byte) {
	addr = (ram.cpu.xint(addr) - ram.baseAddr) / 8
	words := ram.words[addr : addr+int(len(program))/8+1]

	clear(words)
	for i, b := range program {
		shift := (i & 7) * 8
		words[i/8] |= int(b) << shift
	}
}

func (ram *RAM) Access(addr int, data *int, width int, write bool) bool {
	i := (addr - ram.baseAddr) / 8
	if i < 0 || i >= len(ram.words) {
		return false
	}

	if width == 8 {
		if write {
			ram.words[i] = *data
		} else {
			*data = ram.words[i]
		}

		return true
	}

	if shift := (addr & 7) * 8; write {
		mask := -1 << (width * 8)
		ram.words[i] = (ram.words[i] &^ (^mask << shift)) | (*data&^mask)<<shift
	} else {
		*data = int(ram.words[i] >> shift)
	}

	return true
}

func (ram *RAM) Dump(startAddr, byteCount int) {
	i0 := (startAddr - ram.baseAddr) / 8
	i1 := i0 + (byteCount+7)/8
	for i := i0; i < i1; i += 4 {
		fmt.Printf(
			"%016x: %016x %016x %016x %016x\r\n",
			ram.baseAddr+i*8,
			uint(ram.words[i]), uint(ram.words[i+1]),
			uint(ram.words[i+2]), uint(ram.words[i+3]),
		)
	}
}
