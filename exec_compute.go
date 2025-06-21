package rv

func (cpu *CPU) execInstrComputeI(imm, rs1, f3, rd int32) {
	switch f3 {
	case 0b_000: // addi
		cpu.x[rd] = cpu.x[rs1] + imm

	case 0b_010: // slti
		if cpu.x[rs1] < imm {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_011: // sltiu
		if uint32(cpu.x[rs1]) < uint32(imm) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_100: // xori
		cpu.x[rd] = cpu.x[rs1] ^ imm

	case 0b_110: // ori
		cpu.x[rd] = cpu.x[rs1] | imm

	case 0b_111: // andi
		cpu.x[rd] = cpu.x[rs1] & imm

	default:
		shamt := bits(imm, 0, 5)

		switch bits(imm, 5, 7)<<3 | f3 {
		case 0b_0000000_001: // slli
			cpu.x[rd] = cpu.x[rs1] << shamt

		case 0b_0000000_101: // srli
			cpu.x[rd] = int32(uint32(cpu.x[rs1]) >> uint32(shamt))

		case 0b_0100000_101: // srai
			cpu.x[rd] = cpu.x[rs1] >> shamt

		default:
			cpu.instrIllegal = true
		}
	}
}

func (cpu *CPU) execInstrComputeR(f7, rs2, rs1, f3, rd int32) {
	switch f7<<3 | f3 {
	case 0b_0000000_000: // add
		cpu.x[rd] = cpu.x[rs1] + cpu.x[rs2]

	case 0b_0100000_000: // sub
		cpu.x[rd] = cpu.x[rs1] - cpu.x[rs2]

	case 0b_0000000_001: // sll
		shamt := bits(cpu.x[rs2], 0, 5)
		cpu.x[rd] = cpu.x[rs1] << shamt

	case 0b_0000000_010: // slt
		if cpu.x[rs1] < cpu.x[rs2] {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_0000000_011: // sltu
		if uint32(cpu.x[rs1]) < uint32(cpu.x[rs2]) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_0000000_100: // xor
		cpu.x[rd] = cpu.x[rs1] ^ cpu.x[rs2]

	case 0b_0000000_101: // srl
		shamt := bits(cpu.x[rs2], 0, 5)
		cpu.x[rd] = int32(uint32(cpu.x[rs1]) >> uint32(shamt))

	case 0b_0100000_101: // sra
		shamt := bits(cpu.x[rs2], 0, 5)
		cpu.x[rd] = cpu.x[rs1] >> shamt

	case 0b_0000000_110: // or
		cpu.x[rd] = cpu.x[rs1] | cpu.x[rs2]

	case 0b_0000000_111: // and
		cpu.x[rd] = cpu.x[rs1] & cpu.x[rs2]

	default:
		cpu.instrIllegal = true
	}
}
