package instr

import (
	"github.com/temnok/rv/state"
)

var computesM = []func(cpu *state.CPU, op Op){
	0: mul,
	1: mulh,
	2: mulhsu,
	3: mulhu,
	4: div,
	5: divu,
	6: rem,
	7: remu,
}

func computeM(cpu *state.CPU, op Op) {
	computesM[op.F3()](cpu, op)
}
