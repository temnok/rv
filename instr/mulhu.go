package instr

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func Mulhu(cpu *state.State, rd, rs1, rs2 int) {
	var c int

	if cpu.Xlen64() {
		a := uint64(cpu.X[rs1])
		b := uint64(cpu.X[rs2])

		hi, _ := bits.Mul64(a, b)
		c = int(hi)
	} else {
		a := int(uint32(cpu.X[rs1]))
		b := int(uint32(cpu.X[rs2]))

		c = (a * b) >> 32
	}

	cpu.Xset(rd, c)
}
