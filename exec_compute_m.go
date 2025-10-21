package rv

import (
	"github.com/temnok/rv/instr"
)

func (cpu *CPU) execComputeM(rs2, rs1, f3, rd int) {
	switch f3 {
	case 0b_000:
		instr.Mul(&cpu.State, rd, rs1, rs2)
	case 0b_001:
		instr.Mulh(&cpu.State, rd, rs1, rs2)
	case 0b_010:
		instr.Mulhsu(&cpu.State, rd, rs1, rs2)
	case 0b_011:
		instr.Mulhu(&cpu.State, rd, rs1, rs2)
	case 0b_100:
		instr.Div(&cpu.State, rd, rs1, rs2)
	case 0b_101:
		instr.Divu(&cpu.State, rd, rs1, rs2)
	case 0b_110:
		instr.Rem(&cpu.State, rd, rs1, rs2)
	case 0b_111:
		instr.Remu(&cpu.State, rd, rs1, rs2)
	}
}
