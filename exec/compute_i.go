package exec

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var iInstr = []func(cpu *state.CPU, op instr.Op){
	0: instr.Addi,
	1: slli,
	2: instr.Slti,
	3: instr.Sltiu,
	4: instr.Xori,
	5: sr_i,
	6: instr.Ori,
	7: instr.Andi,
}

func ComputeI(cpu *state.CPU, op instr.Op) {
	iInstr[op.F3()](cpu, op)
}

func slli(cpu *state.CPU, op instr.Op) {
	imm := imm.I(op.Code())

	switch imm &^ cpu.Xmask() {
	case 0:
		instr.Slli(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}

func sr_i(cpu *state.CPU, op instr.Op) {
	imm := imm.I(op.Code())

	switch imm &^ cpu.Xmask() {
	case 0:
		instr.Srli(cpu, op)
	case 0b_010000000000:
		instr.Srai(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
