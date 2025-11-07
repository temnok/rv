package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
)

var branch = []func(*state.CPU, instr.Op){
	0: instr.Beq,
	1: instr.Bne,
	2: instrIllegal,
	3: instrIllegal,
	4: instr.Blt,
	5: instr.Bge,
	6: instr.Bltu,
	7: instr.Bgeu,
}

func Branch(cpu *state.CPU, op instr.Op) {
	branch[op.F3()](cpu, op)
}
