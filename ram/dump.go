package ram

import "fmt"

func Dump(ram []int, startAddr, byteCount int) {
	i0 := (startAddr - BaseAddr) / 8
	i1 := i0 + (byteCount+7)/8
	for i := i0; i < i1; i += 4 {
		fmt.Printf(
			"%016x: %016x %016x %016x %016x\r\n",
			BaseAddr+i*8,
			uint(ram[i]), uint(ram[i+1]),
			uint(ram[i+2]), uint(ram[i+3]),
		)
	}
}
