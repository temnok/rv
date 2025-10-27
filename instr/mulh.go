package instr

import (
	"github.com/temnok/rv/state"
	bi "math/bits"
)

func Mulh(cpu *state.State, rd, rs1, rs2 int) {
	a := cpu.X[rs1]
	b := cpu.X[rs2]

	var c int

	if cpu.Xlen64() {
		hi, _ := bi.Mul64(uint64(a), uint64(b))
		s1 := (a >> 63) & b
		s2 := (b >> 63) & a
		c = int(hi) - s1 - s2
	} else {
		c = a * b >> 32
	}

	cpu.Xset(rd, c)
}
