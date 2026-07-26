package ram

import (
	"encoding/binary"
	"io"
	"os"
)

const BaseAddr = 0x8000_0000

func Populate(ram []int, addr int, data []byte) {
	words := ram[(addr-BaseAddr)/8:]

	for i, b := range data {
		shift := (i & 7) * 8
		words[i/8] |= int(b) << shift
	}
}

func PopulateFromFile(ram []int, addr int, path string) {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer r.Close()

	populateFromReader(ram, addr, r)
}

func populateFromReader(ram []int, addr int, r io.ReadCloser) {
	words := ram[(addr-BaseAddr)/8:]
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
