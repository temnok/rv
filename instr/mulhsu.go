package instr

import (
	"github.com/temnok/rv/state"
	bi "math/bits"
)

func Mulhsu(cpu *state.State, rd, rs1, rs2 int) {
	a, b, c := cpu.X[rs1], cpu.X[rs2], 0

	if cpu.XLen64() {
		hi, _ := bi.Mul64(uint64(a), uint64(b))
		s := (a >> 63) & b
		c = int(hi) - s
	} else {
		c = int(int64(a) * int64(uint32(b)) >> 32)
	}

	cpu.Xset(rd, c)
}
