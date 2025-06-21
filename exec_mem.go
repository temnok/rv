package rv

func (cpu *CPU) execInstrLoad(imm, rs1, f3, rd int32) {
	switch f3 {
	case 0b_000: // lb
		if val := cpu.memRead(cpu.x[rs1]+imm, 1); !cpu.faulted() {
			cpu.x[rd] = int32(int8(val))
		}

	case 0b_001: // lh
		if val := cpu.memRead(cpu.x[rs1]+imm, 2); !cpu.faulted() {
			cpu.x[rd] = int32(int16(val))
		}

	case 0b_010: // lw
		if val := cpu.memRead(cpu.x[rs1]+imm, 4); !cpu.faulted() {
			cpu.x[rd] = val
		}

	case 0b_100: // lbu
		if val := cpu.memRead(cpu.x[rs1]+imm, 1); !cpu.faulted() {
			cpu.x[rd] = int32(uint8(val))
		}

	case 0b_101: // lhu
		if val := cpu.memRead(cpu.x[rs1]+imm, 2); !cpu.faulted() {
			cpu.x[rd] = int32(uint16(val))
		}

	default:
		cpu.instrIllegal = true
	}
}

func (cpu *CPU) execInstrStore(imm, rs2, rs1, f3 int32) {
	switch f3 {
	case 0b_000: // sb
		cpu.memWrite(cpu.x[rs1]+imm, 1, cpu.x[rs2])

	case 0b_001: // sh
		cpu.memWrite(cpu.x[rs1]+imm, 2, cpu.x[rs2])

	case 0b_010: // sw
		cpu.memWrite(cpu.x[rs1]+imm, 4, cpu.x[rs2])

	default:
		cpu.instrIllegal = true
	}
}

func (cpu *CPU) execInstrFence(imm, rs1, f3, rd int32) {
	switch f3 {
	case 0b_000: // fence
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			cpu.instrIllegal = true
		}

	case 0b_001: // fence.i
		if imm != 0 || rs1 != 0 || rd != 0 {
			cpu.instrIllegal = true
		}

	default:
		cpu.instrIllegal = true
	}
}
