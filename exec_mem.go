package rv

func (cpu *CPU) execLoad(imm, rs1, f3, rd int) {
	var val int

	switch f3 {
	case 0b_000: // lb
		cpu.memRead(cpu.X[rs1]+imm, &val, 1)
		cpu.xset(rd, int(int8(val)))

	case 0b_001: // lh
		cpu.memRead(cpu.X[rs1]+imm, &val, 2)
		cpu.xset(rd, int(int16(val)))

	case 0b_010: // lw
		cpu.memRead(cpu.X[rs1]+imm, &val, 4)
		cpu.xset(rd, int(int32(val)))

	case 0b_011: // ld
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memRead(cpu.X[rs1]+imm, &val, 8)
		cpu.xset(rd, val)

	case 0b_100: // lbu
		cpu.memRead(cpu.X[rs1]+imm, &val, 1)
		cpu.xset(rd, int(uint8(val)))

	case 0b_101: // lhu
		cpu.memRead(cpu.X[rs1]+imm, &val, 2)
		cpu.xset(rd, int(uint16(val)))

	case 0b_110: // lwu
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memRead(cpu.X[rs1]+imm, &val, 4)
		cpu.xset(rd, int(uint32(val)))
	}

	if cpu.Updated.XReg < 0 && !cpu.isTrapped() {
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execStore(imm, rs2, rs1, f3 int) {
	switch f3 {
	case 0b_000: // sb
		cpu.memWrite(cpu.X[rs1]+imm, int(uint8(cpu.X[rs2])), 1)

	case 0b_001: // sh
		cpu.memWrite(cpu.X[rs1]+imm, int(uint16(cpu.X[rs2])), 2)

	case 0b_010: // sw
		cpu.memWrite(cpu.X[rs1]+imm, int(uint32(cpu.X[rs2])), 4)

	case 0b_011: // sd
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.X[rs1]+imm, cpu.X[rs2], 8)

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execFence(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000: // fence
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			cpu.trap(ExceptionIllegalIstruction)
		}

	case 0b_001: // fence.i
		cpu.Updated.ICache.Clear()

		if imm != 0 || rs1 != 0 || rd != 0 {
			cpu.trap(ExceptionIllegalIstruction)
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
