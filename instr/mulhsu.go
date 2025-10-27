package instr

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func Mulhsu(cpu *state.State, rd, rs1, rs2 int) {
	var c int

	if cpu.XLen64() {
		a := cpu.X[rs1]
		b := cpu.X[rs2]

		hi, _ := bits.Mul64(uint64(a), uint64(b))
		s := (a >> 63) & b
		c = int(hi) - s
	} else {
		a := cpu.X[rs1]
		b := int(uint32(cpu.X[rs2]))

		c = (a * b) >> 32
	}

	cpu.Xset(rd, c)
}
