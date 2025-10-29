package rv

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
)

var mInstr = []func(cpu *state.State, rd, rs1, rs2 int){
	0: instr.Mul,
	1: instr.Mulh,
	2: instr.Mulhsu,
	3: instr.Mulhu,
	4: instr.Div,
	5: instr.Divu,
	6: instr.Rem,
	7: instr.Remu,
}

func (cpu *CPU) execComputeM(rs2, rs1, f3, rd int) {
	mInstr[f3](cpu.State, rd, rs1, rs2)
}
