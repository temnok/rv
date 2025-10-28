package rv

import "github.com/temnok/rv/instr"

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
		cpu.Trap(ExceptionIllegalIstruction)
	}
}
