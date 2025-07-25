package rv

func (cpu *CPU) execComputeI64(imm, rs1, f3, rd int) {
	if cpu.xlen64() {
		switch f3 {
		case 0b_000: // addiw
			cpu.xset(rd, int(int32(cpu.X[rs1])+int32(imm)))

		case 0b_001: // slliw
			if imm < 32 {
				cpu.xset(rd, int(int32(cpu.X[rs1])<<int32(imm)))
			}

		case 0b_101:
			if imm < 32 { // srliw
				cpu.xset(rd, int(int32(uint32(cpu.X[rs1])>>uint32(imm))))
			} else if imm &^= 0b0100000_00000; imm < 32 { // sraiw
				cpu.xset(rd, int(int32(cpu.X[rs1])>>int32(imm)))
			}
		}
	}

	if cpu.Updated.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execComputeR64(f7, rs2, rs1, f3, rd int) {
	if !cpu.xlen64() {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if f7 == 1 {
		cpu.execComputeM64(rs2, rs1, f3, rd)
		return
	}

	op := bit(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000: // addw
		cpu.xset(rd, int(int32(cpu.X[rs1])+int32(cpu.X[rs2])))

	case 0b_1_000: // subw
		cpu.xset(rd, int(int32(cpu.X[rs1])-int32(cpu.X[rs2])))

	case 0b_001: // sllw
		shamt := int32(cpu.X[rs2]) & 31
		cpu.xset(rd, int(int32(cpu.X[rs1])<<shamt))

	case 0b_101: // srlw
		shamt := uint32(cpu.X[rs2]) & 31
		cpu.xset(rd, int(int32(uint32(cpu.X[rs1])>>shamt)))

	case 0b_1_101: // sraw
		shamt := int32(cpu.X[rs2]) & 31
		cpu.xset(rd, int(int32(cpu.X[rs1])>>shamt))
	}

	if cpu.Updated.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}
