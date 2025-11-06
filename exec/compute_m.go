package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
)

var mInstr = []func(cpu *state.CPU, op instr.Op){
	0: instr.Mul,
	1: instr.Mulh,
	2: instr.Mulhsu,
	3: instr.Mulhu,
	4: instr.Div,
	5: instr.Divu,
	6: instr.Rem,
	7: instr.Remu,
}

func ComputeM(cpu *state.CPU, op instr.Op) {
	mInstr[op.F3()](cpu, op)
}
