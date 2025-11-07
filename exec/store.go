package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
)

var store = []func(*state.CPU, instr.Op){
	0: instr.Sb,
	1: instr.Sh,
	2: instr.Sw,
	3: instr.Sd,
	4: instrIllegal,
	5: instrIllegal,
	6: instrIllegal,
	7: instrIllegal,
}

func Store(cpu *state.CPU, op instr.Op) {
	store[op.F3()](cpu, op)
}
