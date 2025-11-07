package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
)

var load = []func(*state.CPU, instr.Op){
	0: instr.Lb,
	1: instr.Lh,
	2: instr.Lw,
	3: instr.Ld,
	4: instr.Lbu,
	5: instr.Lhu,
	6: instr.Lwu,
	7: instrIllegal,
}

func Load(cpu *state.CPU, op instr.Op) {
	load[op.F3()](cpu, op)
}
