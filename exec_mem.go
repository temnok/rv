package rv

func (cpu *CPU) execLoad(imm, rs1, f3, rd int) {
	var val int

	switch f3 {
	case 0b_000: // lb
		if cpu.memRead(cpu.x[rs1]+imm, &val, 1); !cpu.isTrapped {
			cpu.x[rd] = int(int8(val))
		}

	case 0b_001: // lh
		if cpu.memRead(cpu.x[rs1]+imm, &val, 2); !cpu.isTrapped {
			cpu.x[rd] = int(int16(val))
		}

	case 0b_010: // lw
		if cpu.memRead(cpu.x[rs1]+imm, &val, 4); !cpu.isTrapped {
			cpu.x[rd] = int(int32(val))
		}

	case 0b_011: // ld
		if Xlen32 {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.x[rs1]+imm, &val, 8); !cpu.isTrapped {
			cpu.x[rd] = int(val)
		}

	case 0b_100: // lbu
		if cpu.memRead(cpu.x[rs1]+imm, &val, 1); !cpu.isTrapped {
			cpu.x[rd] = int(uint8(val))
		}

	case 0b_101: // lhu
		if cpu.memRead(cpu.x[rs1]+imm, &val, 2); !cpu.isTrapped {
			cpu.x[rd] = int(uint16(val))
		}

	case 0b_110: // lwu
		if Xlen32 {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.x[rs1]+imm, &val, 4); !cpu.isTrapped {
			cpu.x[rd] = int(uint32(val))
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execStore(imm, rs2, rs1, f3 int) {
	switch f3 {
	case 0b_000: // sb
		cpu.memWrite(cpu.x[rs1]+imm, int(uint8(cpu.x[rs2])), 1)

	case 0b_001: // sh
		cpu.memWrite(cpu.x[rs1]+imm, int(uint16(cpu.x[rs2])), 2)

	case 0b_010: // sw
		cpu.memWrite(cpu.x[rs1]+imm, int(uint32(cpu.x[rs2])), 4)

	case 0b_011: // sd
		if Xlen32 {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.x[rs1]+imm, cpu.x[rs2], 8)

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
		if imm != 0 || rs1 != 0 || rd != 0 {
			cpu.trap(ExceptionIllegalIstruction)
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
