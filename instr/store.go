package instr

import (
	"github.com/temnok/rv/state"
)

var stores = []func(*state.CPU, Op){
	0: sb,
	1: sh,
	2: sw,
	3: sd,
	4: illegal,
	5: illegal,
	6: illegal,
	7: illegal,
}

func Store(cpu *state.CPU, op Op) {
	stores[op.F3()](cpu, op)
}
