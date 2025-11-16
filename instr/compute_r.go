package instr

import (
	"github.com/temnok/rv/state"
)

var computesR = []instr{
	0: add,
	1: sll,
	2: slt,
	3: sltu,
	4: xor,
	5: srl,
	6: or,
	7: and,
}

func ComputeR(cpu *state.CPU, op Op) {
	f3, f7 := op.F3(), op.F7()

	ins := illegal

	switch f7 {
	case 0:
		ins = computesR[f3]
	case 1:
		ins = computeM
	case 32:
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
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	c := f(a, b)

	cpu.Xset(op.Rd(), c)
}
