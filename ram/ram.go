package ram

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type RAM struct {
	words []int
}

const BaseAddr = 0x8000_0000

func New(size int) *RAM {
	return &RAM{
		words: make([]int, size/8),
	}
}

func (ram *RAM) Load(addr int, data []byte) {
	words := ram.words[(addr-BaseAddr)/8:]

	for i, b := range data {
		shift := (i & 7) * 8
		words[i/8] |= int(b) << shift
	}
}

func (ram *RAM) LoadFile(addr int, path string) {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer r.Close()

	ram.LoadReader(addr, r)
}

func (ram *RAM) LoadReader(addr int, r io.ReadCloser) {
	words := ram.words[(addr-BaseAddr)/8:]
	buf := make([]byte, 8*1024)

	for {
		n, _ := r.Read(buf)
		for ; n%8 != 0; n++ {
			buf[n] = 0
		}

		for i := 0; i < n; i += 8 {
			words[0] = int(binary.LittleEndian.Uint64(buf[i:]))
			words = words[1:]
		}

		if n < len(buf) {
			break
		}
	}
}

func (ram *RAM) Access(addr int, width int, write bool, writeData int) int {
	i := (addr - BaseAddr) / 8
	if i < 0 || i >= len(ram.words) {
		return 0
	}

	if width == 8 {
		data := ram.words[i]

		if write {
			ram.words[i] = writeData
		}

		return data
	}

	shift := (addr & 7) * 8
	mask := 1<<(width*8) - 1
	val := ram.words[i]

	if write {
		ram.words[i] = (val &^ (mask << shift)) | (writeData&mask)<<shift
	}

	return val >> shift & mask
}

func (ram *RAM) Dump(startAddr, byteCount int) {
	i0 := (startAddr - BaseAddr) / 8
	i1 := i0 + (byteCount+7)/8
	for i := i0; i < i1; i += 4 {
		fmt.Printf(
			"%016x: %016x %016x %016x %016x\r\n",
			BaseAddr+i*8,
			uint(ram.words[i]), uint(ram.words[i+1]),
			uint(ram.words[i+2]), uint(ram.words[i+3]),
		)
	}
}
