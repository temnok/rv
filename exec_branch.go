package rv

func (cpu *CPU) execBranch(imm, rs2, rs1, f3 int) {
	cond := false

	switch f3 {
	case 0b_000: // beq
		cond = cpu.Reg[rs1] == cpu.Reg[rs2]

	case 0b_001: // bne
		cond = cpu.Reg[rs1] != cpu.Reg[rs2]

	case 0b_100: // blt
		cond = cpu.Reg[rs1] < cpu.Reg[rs2]

	case 0b_101: // bge
		cond = cpu.Reg[rs1] >= cpu.Reg[rs2]

	case 0b_110: // bltu
		cond = cpu.xuint(cpu.Reg[rs1]) < cpu.xuint(cpu.Reg[rs2])

	case 0b_111: // bgeu
		cond = cpu.xuint(cpu.Reg[rs1]) >= cpu.xuint(cpu.Reg[rs2])

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}

	if cond {
		cpu.Updated.PC = cpu.xint(cpu.PC + imm)
	}
}
