package rv

func (cpu *CPU) execBranch(imm, rs2, rs1, f3 int) {
	cond := false

	switch f3 {
	case 0b_000: // beq
		cond = cpu.X[rs1] == cpu.X[rs2]

	case 0b_001: // bne
		cond = cpu.X[rs1] != cpu.X[rs2]

	case 0b_100: // blt
		cond = cpu.X[rs1] < cpu.X[rs2]

	case 0b_101: // bge
		cond = cpu.X[rs1] >= cpu.X[rs2]

	case 0b_110: // bltu
		cond = cpu.Xuint(cpu.X[rs1]) < cpu.Xuint(cpu.X[rs2])

	case 0b_111: // bgeu
		cond = cpu.Xuint(cpu.X[rs1]) >= cpu.Xuint(cpu.X[rs2])

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}

	if cond {
		cpu.Update.PC = cpu.Xint(cpu.PC + imm)
	}
}
