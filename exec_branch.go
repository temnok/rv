package rv

func (cpu *CPU) execInstrBranch(imm, rs2, rs1, f3 int32) {
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
		cond = uint32(cpu.x[rs1]) < uint32(cpu.x[rs2])

	case 0b_111: // bgeu
		cond = uint32(cpu.x[rs1]) >= uint32(cpu.x[rs2])

	default:
		cpu.instrIllegal = true
	}

	if cond {
		cpu.nextPC = cpu.pc + imm
	}
}
