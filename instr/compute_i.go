package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var iInstr = []func(cpu *state.CPU, op Op){
	0: Addi,
	1: slli,
	2: Slti,
	3: Sltiu,
	4: Xori,
	5: sr_i,
	6: Ori,
	7: Andi,
}

func ComputeI(cpu *state.CPU, op Op) {
	iInstr[op.F3()](cpu, op)
}

func slli(cpu *state.CPU, op Op) {
	imm := imm.I(op.Code())

	switch imm &^ cpu.Xmask() {
	case 0:
		Slli(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}

func sr_i(cpu *state.CPU, op Op) {
	imm := imm.I(op.Code())

	switch imm &^ cpu.Xmask() {
	case 0:
		Srli(cpu, op)
	case 0b_010000000000:
		Srai(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
