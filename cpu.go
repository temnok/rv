package rv

const (
	mepc = 0x341
)

type CPU struct {
	x          [32]int32
	pc, nextPC int32
	csr        [4096]int32
	mem        []byte

	instrIllegal   bool
	memAccessFault bool
	eCall          bool
}

func (cpu *CPU) faulted() bool {
	return cpu.instrIllegal || cpu.memAccessFault || cpu.eCall
}
