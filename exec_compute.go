package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.xSet(rd, cpu.X[rs1]+imm)

	case 0b_001: // slli
		if imm < cpu.XLen {
			cpu.xSet(rd, cpu.X[rs1]<<imm)
		}

	case 0b_010: // slti
		if cpu.X[rs1] < imm {
			cpu.xSet(rd, 1)
		} else {
			cpu.xSet(rd, 0)
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(imm) {
			cpu.xSet(rd, 1)
		} else {
			cpu.xSet(rd, 0)
		}

	case 0b_100: // xori
		cpu.xSet(rd, cpu.X[rs1]^imm)

	case 0b_101:
		if imm < cpu.XLen { // srli
			cpu.xSet(rd, int(cpu.xuint(cpu.X[rs1])>>cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
			cpu.xSet(rd, cpu.X[rs1]>>imm)
		}

	case 0b_110: // ori
		cpu.xSet(rd, cpu.X[rs1]|imm)

	case 0b_111: // andi
		cpu.xSet(rd, cpu.X[rs1]&imm)
	}

	if cpu.Updated.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
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
		cpu.xSet(rd, cpu.X[rs1]+cpu.X[rs2])

	case 0b_1_000: // sub
		cpu.xSet(rd, cpu.X[rs1]-cpu.X[rs2])

	case 0b_001: // sll
		cpu.xSet(rd, cpu.X[rs1]<<(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_010: // slt
		if cpu.X[rs1] < cpu.X[rs2] {
			cpu.xSet(rd, 1)
		} else {
			cpu.xSet(rd, 0)
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(cpu.X[rs2]) {
			cpu.xSet(rd, 1)
		} else {
			cpu.xSet(rd, 0)
		}

	case 0b_100: // xor
		cpu.xSet(rd, cpu.X[rs1]^cpu.X[rs2])

	case 0b_101: // srl
		cpu.xSet(rd, int(cpu.xuint(cpu.X[rs1])>>cpu.xuint(cpu.X[rs2]&(cpu.XLen-1))))

	case 0b_1_101: // sra
		cpu.xSet(rd, cpu.X[rs1]>>(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_110: // or
		cpu.xSet(rd, cpu.X[rs1]|cpu.X[rs2])

	case 0b_111: // and
		cpu.xSet(rd, cpu.X[rs1]&cpu.X[rs2])
	}

	if cpu.Updated.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}
