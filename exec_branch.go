package rv

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execBranch(imm, rs2, rs1, f3 int) {
	switch f3 {
	case 0b_000:
		instr.Beq(cpu.State, rs1, rs2, imm)
	case 0b_001:
		instr.Bne(cpu.State, rs1, rs2, imm)
	case 0b_100:
		instr.Blt(cpu.State, rs1, rs2, imm)
	case 0b_101:
		instr.Bge(cpu.State, rs1, rs2, imm)
	case 0b_110:
		instr.Bltu(cpu.State, rs1, rs2, imm)
	case 0b_111:
		instr.Bgeu(cpu.State, rs1, rs2, imm)
	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
