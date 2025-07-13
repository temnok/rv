package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.Updated.RegValue = cpu.xint(cpu.Reg[rs1] + imm)

	case 0b_001: // slli
		if imm < cpu.XLen {
			cpu.Updated.RegValue = cpu.xint(cpu.Reg[rs1] << imm)
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_010: // slti
		if cpu.Reg[rs1] < imm {
			cpu.Updated.RegValue = 1
		} else {
			cpu.Updated.RegValue = 0
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.Reg[rs1]) < cpu.xuint(imm) {
			cpu.Updated.RegValue = 1
		} else {
			cpu.Updated.RegValue = 0
		}

	case 0b_100: // xori
		cpu.Updated.RegValue = cpu.Reg[rs1] ^ imm

	case 0b_101:
		if imm < cpu.XLen { // srli
			cpu.Updated.RegValue = cpu.xint(int(cpu.xuint(cpu.Reg[rs1]) >> cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
			cpu.Updated.RegValue = cpu.Reg[rs1] >> imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_110: // ori
		cpu.Updated.RegValue = cpu.Reg[rs1] | imm

	case 0b_111: // andi
		cpu.Updated.RegValue = cpu.Reg[rs1] & imm
	}

	if !cpu.isTrapped() {
		cpu.Updated.RegIndex = rd
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
		cpu.Updated.RegValue = cpu.xint(cpu.Reg[rs1] + cpu.Reg[rs2])

	case 0b_1_000: // sub
		cpu.Updated.RegValue = cpu.xint(cpu.Reg[rs1] - cpu.Reg[rs2])

	case 0b_001: // sll
		cpu.Updated.RegValue = cpu.xint(cpu.Reg[rs1] << (cpu.Reg[rs2] & (cpu.XLen - 1)))

	case 0b_010: // slt
		if cpu.Reg[rs1] < cpu.Reg[rs2] {
			cpu.Updated.RegValue = 1
		} else {
			cpu.Updated.RegValue = 0
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.Reg[rs1]) < cpu.xuint(cpu.Reg[rs2]) {
			cpu.Updated.RegValue = 1
		} else {
			cpu.Updated.RegValue = 0
		}

	case 0b_100: // xor
		cpu.Updated.RegValue = cpu.Reg[rs1] ^ cpu.Reg[rs2]

	case 0b_101: // srl
		cpu.Updated.RegValue = cpu.xint(int(cpu.xuint(cpu.Reg[rs1]) >> cpu.xuint(cpu.Reg[rs2]&(cpu.XLen-1))))

	case 0b_1_101: // sra
		cpu.Updated.RegValue = cpu.Reg[rs1] >> (cpu.Reg[rs2] & (cpu.XLen - 1))

	case 0b_110: // or
		cpu.Updated.RegValue = cpu.Reg[rs1] | cpu.Reg[rs2]

	case 0b_111: // and
		cpu.Updated.RegValue = cpu.Reg[rs1] & cpu.Reg[rs2]

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.RegIndex = rd
}
