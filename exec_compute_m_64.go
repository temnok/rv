package rv

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execComputeM64(rs2, rs1, f3, rd int) {
	switch f3 {
	case 0b_000:
		instr.Mulw(&cpu.State, rd, rs1, rs2)
	case 0b_100:
		instr.Divw(&cpu.State, rd, rs1, rs2)
	case 0b_101:
		instr.Divuw(&cpu.State, rd, rs1, rs2)
	case 0b_110:
		instr.Remw(&cpu.State, rd, rs1, rs2)
	case 0b_111:
		instr.Remuw(&cpu.State, rd, rs1, rs2)
	default:
		trap.EnterWithoutTval(&cpu.State, ExceptionIllegalIstruction)
	}
}
