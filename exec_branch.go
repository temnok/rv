package rv

func (cpu *CPU) execBranch(imm, rs2, rs1, f3 int) {
	cond := false

	switch f3 {
	case 0b_000: // beq
		cond = cpu.x[rs1] == cpu.x[rs2]

	case 0b_001: // bne
		cond = cpu.x[rs1] != cpu.x[rs2]

	case 0b_100: // blt
		cond = cpu.x[rs1] < cpu.x[rs2]

	case 0b_101: // bge
		cond = cpu.x[rs1] >= cpu.x[rs2]

	case 0b_110: // bltu
		cond = cpu.xuint(cpu.x[rs1]) < cpu.xuint(cpu.x[rs2])

	case 0b_111: // bgeu
		cond = cpu.xuint(cpu.x[rs1]) >= cpu.xuint(cpu.x[rs2])

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}

	if cond {
		cpu.next.pc = cpu.xint(cpu.pc + imm)
	}
}
