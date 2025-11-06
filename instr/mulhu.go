package instr

import (
	"github.com/temnok/rv/state"
	"math/bits"
)

func Mulhu(cpu *state.CPU, op Op) {
	var c int

	if cpu.Xlen64() {
		a := uint64(cpu.X[op.Rs1()])
		b := uint64(cpu.X[op.Rs2()])

		hi, _ := bits.Mul64(a, b)
		c = int(hi)
	} else {
		a := int(uint32(cpu.X[op.Rs1()]))
		b := int(uint32(cpu.X[op.Rs2()]))

		c = (a * b) >> 32
	}

	cpu.Xset(op.Rd(), c)
}
