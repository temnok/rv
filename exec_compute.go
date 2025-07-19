package rv

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // addi
		cpu.xset(rd, cpu.X[rs1]+imm)

	case 0b_001: // slli
		if imm < cpu.XLen {
			cpu.xset(rd, cpu.X[rs1]<<imm)
		}

	case 0b_010: // slti
		if cpu.X[rs1] < imm {
			cpu.xset(rd, 1)
		} else {
			cpu.xset(rd, 0)
		}

	case 0b_011: // sltiu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(imm) {
			cpu.xset(rd, 1)
		} else {
			cpu.xset(rd, 0)
		}

	case 0b_100: // xori
		cpu.xset(rd, cpu.X[rs1]^imm)

	case 0b_101:
		if imm < cpu.XLen { // srli
			cpu.xset(rd, int(cpu.xuint(cpu.X[rs1])>>cpu.xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
			cpu.xset(rd, cpu.X[rs1]>>imm)
		}

	case 0b_110: // ori
		cpu.xset(rd, cpu.X[rs1]|imm)

	case 0b_111: // andi
		cpu.xset(rd, cpu.X[rs1]&imm)
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
		cpu.xset(rd, cpu.X[rs1]+cpu.X[rs2])

	case 0b_1_000: // sub
		cpu.xset(rd, cpu.X[rs1]-cpu.X[rs2])

	case 0b_001: // sll
		cpu.xset(rd, cpu.X[rs1]<<(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_010: // slt
		if cpu.X[rs1] < cpu.X[rs2] {
			cpu.xset(rd, 1)
		} else {
			cpu.xset(rd, 0)
		}

	case 0b_011: // sltu
		if cpu.xuint(cpu.X[rs1]) < cpu.xuint(cpu.X[rs2]) {
			cpu.xset(rd, 1)
		} else {
			cpu.xset(rd, 0)
		}

	case 0b_100: // xor
		cpu.xset(rd, cpu.X[rs1]^cpu.X[rs2])

	case 0b_101: // srl
		cpu.xset(rd, int(cpu.xuint(cpu.X[rs1])>>cpu.xuint(cpu.X[rs2]&(cpu.XLen-1))))

	case 0b_1_101: // sra
		cpu.xset(rd, cpu.X[rs1]>>(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_110: // or
		cpu.xset(rd, cpu.X[rs1]|cpu.X[rs2])

	case 0b_111: // and
		cpu.xset(rd, cpu.X[rs1]&cpu.X[rs2])
	}

	if cpu.Updated.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}
