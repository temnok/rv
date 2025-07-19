package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.Updated.XVal = cpu.xint(cpu.X[rs1] + imm)

	case 0b_001: // slli
		if imm < cpu.XLen {
			cpu.Updated.XVal = cpu.xint(cpu.X[rs1] << imm)
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_010: // slti
		if cpu.X[rs1] < imm {
			cpu.Updated.XVal = 1
		} else {
			cpu.Updated.XVal = 0
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(imm) {
			cpu.Updated.XVal = 1
		} else {
			cpu.Updated.XVal = 0
		}

	case 0b_100: // xori
		cpu.Updated.XVal = cpu.X[rs1] ^ imm

	case 0b_101:
		if imm < cpu.XLen { // srli
			cpu.Updated.XVal = cpu.xint(int(cpu.xuint(cpu.X[rs1]) >> cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
			cpu.Updated.XVal = cpu.X[rs1] >> imm
		} else {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_110: // ori
		cpu.Updated.XVal = cpu.X[rs1] | imm

	case 0b_111: // andi
		cpu.Updated.XVal = cpu.X[rs1] & imm
	}

	if !cpu.isTrapped() {
		cpu.Updated.XReg = rd
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
		cpu.Updated.XVal = cpu.xint(cpu.X[rs1] + cpu.X[rs2])

	case 0b_1_000: // sub
		cpu.Updated.XVal = cpu.xint(cpu.X[rs1] - cpu.X[rs2])

	case 0b_001: // sll
		cpu.Updated.XVal = cpu.xint(cpu.X[rs1] << (cpu.X[rs2] & (cpu.XLen - 1)))

	case 0b_010: // slt
		if cpu.X[rs1] < cpu.X[rs2] {
			cpu.Updated.XVal = 1
		} else {
			cpu.Updated.XVal = 0
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(cpu.X[rs2]) {
			cpu.Updated.XVal = 1
		} else {
			cpu.Updated.XVal = 0
		}

	case 0b_100: // xor
		cpu.Updated.XVal = cpu.X[rs1] ^ cpu.X[rs2]

	case 0b_101: // srl
		cpu.Updated.XVal = cpu.xint(int(cpu.xuint(cpu.X[rs1]) >> cpu.xuint(cpu.X[rs2]&(cpu.XLen-1))))

	case 0b_1_101: // sra
		cpu.Updated.XVal = cpu.X[rs1] >> (cpu.X[rs2] & (cpu.XLen - 1))

	case 0b_110: // or
		cpu.Updated.XVal = cpu.X[rs1] | cpu.X[rs2]

	case 0b_111: // and
		cpu.Updated.XVal = cpu.X[rs1] & cpu.X[rs2]

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.XReg = rd
}
