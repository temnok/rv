package instr

import (
	"github.com/temnok/rv/state"
)

var insComputeM64 = []instr{
	0: mulw,
	1: illegal,
	2: illegal,
	3: illegal,
	4: divw,
	5: divuw,
	6: remw,
	7: remuw,
}

func computeM64(cpu *state.CPU, op Op) {
	insComputeM64[op.F3()](cpu, op)
}
