package rv

func (cpu *CPU) execComputeIw(imm, rs1, f3, rd int) {
	if !cpu.xlen64() {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch f3 {
	case 0b_000: // addiw
		cpu.updated.regIndex = rd
		cpu.updated.regValue = int(int32(cpu.reg[rs1]) + int32(imm))
		return

	case 0b_001: // slliw
		if imm < 32 {
			cpu.updated.regIndex = rd
			cpu.updated.regValue = int(int32(cpu.reg[rs1]) << int32(imm))
			return
		}

	case 0b_101:
		if imm < 32 { // srliw
			cpu.updated.regIndex = rd
			cpu.updated.regValue = int(int32(uint32(cpu.reg[rs1]) >> uint32(imm)))
			return
		} else if imm &^= 0b0100000_00000; imm < 32 { // sraiw
			cpu.updated.regIndex = rd
			cpu.updated.regValue = int(int32(cpu.reg[rs1]) >> int32(imm))
			return
		}
	}

	cpu.trap(ExceptionIllegalIstruction)
}

func (cpu *CPU) execComputeRw(f7, rs2, rs1, f3, rd int) {
	if !cpu.xlen64() {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if f7 == 1 {
		cpu.execComputeMw(rs2, rs1, f3, rd)
		return
	}

	op := bit(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000: // addw
		cpu.updated.regValue = int(int32(cpu.reg[rs1]) + int32(cpu.reg[rs2]))

	case 0b_1_000: // subw
		cpu.updated.regValue = int(int32(cpu.reg[rs1]) - int32(cpu.reg[rs2]))

	case 0b_001: // sllw
		shamt := int32(cpu.reg[rs2]) & 31
		cpu.updated.regValue = int(int32(cpu.reg[rs1]) << shamt)

	case 0b_101: // srlw
		shamt := uint32(cpu.reg[rs2]) & 31
		cpu.updated.regValue = int(int32(uint32(cpu.reg[rs1]) >> shamt))

	case 0b_1_101: // sraw
		shamt := int32(cpu.reg[rs2]) & 31
		cpu.updated.regValue = int(int32(cpu.reg[rs1]) >> shamt)

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.updated.regIndex = rd
}
