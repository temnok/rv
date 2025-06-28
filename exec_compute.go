package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd Xint) {
	switch f3 {
	case 0b_000: // addi
		cpu.x[rd] = cpu.x[rs1] + imm

	case 0b_001: // slli
		if imm < Xbits {
			cpu.x[rd] = cpu.x[rs1] << imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_010: // slti
		if cpu.x[rs1] < imm {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_011: // sltiu
		if Xuint(cpu.x[rs1]) < Xuint(imm) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_100: // xori
		cpu.x[rd] = cpu.x[rs1] ^ imm

	case 0b_101:
		if imm < Xbits { // srli
			cpu.x[rd] = Xint(Xuint(cpu.x[rs1]) >> Xuint(imm))
		} else if imm &^= 0b0100000_00000; imm < Xbits { // srai
			cpu.x[rd] = cpu.x[rs1] >> imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_110: // ori
		cpu.x[rd] = cpu.x[rs1] | imm

	case 0b_111: // andi
		cpu.x[rd] = cpu.x[rs1] & imm
	}
}

func (cpu *CPU) execComputeR(f7, rs2, rs1, f3, rd Xint) {
	if f7 == 1 {
		cpu.execComputeM(rs2, rs1, f3, rd)
		return
	}

	switch f7<<3 | f3 {
	case 0b_0000000_000: // add
		cpu.x[rd] = cpu.x[rs1] + cpu.x[rs2]

	case 0b_0100000_000: // sub
		cpu.x[rd] = cpu.x[rs1] - cpu.x[rs2]

	case 0b_0000000_001: // sll
		shamt := bits(cpu.x[rs2], 0, Xshift)
		cpu.x[rd] = cpu.x[rs1] << shamt

	case 0b_0000000_010: // slt
		if cpu.x[rs1] < cpu.x[rs2] {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_0000000_011: // sltu
		if Xuint(cpu.x[rs1]) < Xuint(cpu.x[rs2]) {
			cpu.x[rd] = 1
		} else {
			cpu.x[rd] = 0
		}

	case 0b_0000000_100: // xor
		cpu.x[rd] = cpu.x[rs1] ^ cpu.x[rs2]

	case 0b_0000000_101: // srl
		shamt := bits(cpu.x[rs2], 0, Xshift)
		cpu.x[rd] = Xint(Xuint(cpu.x[rs1]) >> Xuint(shamt))

	case 0b_0100000_101: // sra
		shamt := bits(cpu.x[rs2], 0, Xshift)
		cpu.x[rd] = cpu.x[rs1] >> shamt

	case 0b_0000000_110: // or
		cpu.x[rd] = cpu.x[rs1] | cpu.x[rs2]

	case 0b_0000000_111: // and
		cpu.x[rd] = cpu.x[rs1] & cpu.x[rs2]

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
