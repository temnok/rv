package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

var branchIns = []func(*state.CPU, Op){
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
	branchIns[op.f3()](cpu, op)
}

func branch(cpu *state.CPU, op Op, f func(a, b int) bool) {
	a := cpu.X[op.rs1()]
	b := cpu.X[op.rs2()]

	if f(a, b) {
		cpu.Update.PC = cpu.Int(cpu.PC + imm.B(op.code()))
	}
}
