package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.updated.regValue = cpu.xint(cpu.reg[rs1] + imm)

	case 0b_001: // slli
		if imm < cpu.xlen {
			cpu.updated.regValue = cpu.xint(cpu.reg[rs1] << imm)
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_010: // slti
		if cpu.reg[rs1] < imm {
			cpu.updated.regValue = 1
		} else {
			cpu.updated.regValue = 0
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.reg[rs1]) < cpu.xuint(imm) {
			cpu.updated.regValue = 1
		} else {
			cpu.updated.regValue = 0
		}

	case 0b_100: // xori
		cpu.updated.regValue = cpu.reg[rs1] ^ imm

	case 0b_101:
		if imm < cpu.xlen { // srli
			cpu.updated.regValue = cpu.xint(int(cpu.xuint(cpu.reg[rs1]) >> cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.xlen { // srai
			cpu.updated.regValue = cpu.reg[rs1] >> imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_110: // ori
		cpu.updated.regValue = cpu.reg[rs1] | imm

	case 0b_111: // andi
		cpu.updated.regValue = cpu.reg[rs1] & imm
	}

	if !cpu.isTrapped {
		cpu.updated.regIndex = rd
	}
}

func (cpu *CPU) execComputeR(f7, rs2, rs1, f3, rd int) {
	if f7 == 1 {
		cpu.execComputeM(rs2, rs1, f3, rd)
		return
	}

	op := bit(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000: // add
		cpu.updated.regValue = cpu.xint(cpu.reg[rs1] + cpu.reg[rs2])

	case 0b_1_000: // sub
		cpu.updated.regValue = cpu.xint(cpu.reg[rs1] - cpu.reg[rs2])

	case 0b_001: // sll
		cpu.updated.regValue = cpu.xint(cpu.reg[rs1] << (cpu.reg[rs2] & (cpu.xlen - 1)))

	case 0b_010: // slt
		if cpu.reg[rs1] < cpu.reg[rs2] {
			cpu.updated.regValue = 1
		} else {
			cpu.updated.regValue = 0
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.reg[rs1]) < cpu.xuint(cpu.reg[rs2]) {
			cpu.updated.regValue = 1
		} else {
			cpu.updated.regValue = 0
		}

	case 0b_100: // xor
		cpu.updated.regValue = cpu.reg[rs1] ^ cpu.reg[rs2]

	case 0b_101: // srl
		cpu.updated.regValue = cpu.xint(int(cpu.xuint(cpu.reg[rs1]) >> cpu.xuint(cpu.reg[rs2]&(cpu.xlen-1))))

	case 0b_1_101: // sra
		cpu.updated.regValue = cpu.reg[rs1] >> (cpu.reg[rs2] & (cpu.xlen - 1))

	case 0b_110: // or
		cpu.updated.regValue = cpu.reg[rs1] | cpu.reg[rs2]

	case 0b_111: // and
		cpu.updated.regValue = cpu.reg[rs1] & cpu.reg[rs2]

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.updated.regIndex = rd
}
