package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

var branches = []func(*state.CPU, Op){
	0: beq,
	1: bne,
	2: illegal,
	3: illegal,
	4: blt,
	5: bge,
	6: bltu,
	7: bgeu,
}

func execBranch(cpu *state.CPU, op Op) {
	branches[op.F3()](cpu, op)
}

func branch(cpu *state.CPU, op Op, f func(a, b int) bool) {
	a := cpu.X[op.Rs1()]
	b := cpu.X[op.Rs2()]

	if f(a, b) {
		cpu.Update.PC = cpu.Int(cpu.PC + imm.B(op.Code()))
	}
}
