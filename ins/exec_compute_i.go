package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

var insComputeI = []func(*state.CPU, Op){
	0: addi,
	1: slli_illegal,
	2: slti,
	3: sltiu,
	4: xori,
	5: srli_srai,
	6: ori,
	7: andi,
}

func execComputeI(cpu *state.CPU, op Op) {
	insComputeI[op.f3()](cpu, op)
}

func computeI(cpu *state.CPU, op Op, f func(a, b int) int) {
	a := cpu.X[op.rs1()]
	b := imm.I(op.code())

	c := f(a, b)

	cpu.Xset(op.rd(), c)
}

func slli_illegal(cpu *state.CPU, op Op) {
	ins := illegal

	if imm.I(op.code())&^cpu.Mask() == 0 {
		ins = slli
	}

	ins(cpu, op)
}

func srli_srai(cpu *state.CPU, op Op) {
	ins := illegal

	switch imm.I(op.code()) &^ cpu.Mask() {
	case 0:
		ins = srli
	case 1 << 10:
		ins = srai
	}

	ins(cpu, op)
}
