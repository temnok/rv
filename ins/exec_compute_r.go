package ins

import (
	"github.com/temnok/rv/state"
)

var computeRIns = []func(*state.CPU, Op){
	0: add,
	1: sll,
	2: slt,
	3: sltu,
	4: xor,
	5: srl,
	6: or,
	7: and,
}

func execComputeR(cpu *state.CPU, op Op) {
	f3, f7 := op.f3(), op.f7()

	ins := illegal

	switch f7 {
	case 0:
		ins = computeRIns[f3]
	case 1:
		ins = computeM
	case 1 << 5:
		switch f3 {
		case 0:
			ins = sub
		case 5:
			ins = sra
		}
	}

	ins(cpu, op)
}

func computeR(cpu *state.CPU, op Op, f func(a, b int) int) {
	a := cpu.X[op.rs1()]
	b := cpu.X[op.rs2()]

	c := f(a, b)

	cpu.Xset(op.rd(), c)
}
