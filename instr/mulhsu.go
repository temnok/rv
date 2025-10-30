package instr

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func Mulhsu(cpu *state.State, op Op) {
	var c int

	if cpu.Xlen64() {
		a := cpu.X[op.Rs1()]
		b := cpu.X[op.Rs2()]

		hi, _ := bits.Mul64(uint64(a), uint64(b))
		s := (a >> 63) & b
		c = int(hi) - s
	} else {
		a := cpu.X[op.Rs1()]
		b := int(uint32(cpu.X[op.Rs2()]))

		c = (a * b) >> 32
	}

	cpu.Xset(op.Rd(), c)
}
