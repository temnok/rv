package rv

func (cpu *CPU) execInstrLoad(imm, rs1, f3, rd int32) {
	success := false
	var val int32

	switch f3 {
	case 0b_000: // lb
		if success = cpu.memRead(cpu.x[rs1]+imm, 1, &val); success {
			cpu.x[rd] = int32(int8(val))
		}

	case 0b_001: // lh
		if success = cpu.memRead(cpu.x[rs1]+imm, 2, &val); success {
			cpu.x[rd] = int32(int16(val))
		}

	case 0b_010: // lw
		if success = cpu.memRead(cpu.x[rs1]+imm, 4, &val); success {
			cpu.x[rd] = val
		}

	case 0b_100: // lbu
		if success = cpu.memRead(cpu.x[rs1]+imm, 1, &val); success {
			cpu.x[rd] = int32(uint8(val))
		}

	case 0b_101: // lhu
		if success = cpu.memRead(cpu.x[rs1]+imm, 2, &val); success {
			cpu.x[rd] = int32(uint16(val))
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if !success {
		cpu.trap(ExceptionLoadAccessFault)
	}
}

func (cpu *CPU) execInstrStore(imm, rs2, rs1, f3 int32) {
	success := false

	switch f3 {
	case 0b_000: // sb
		success = cpu.memWrite(cpu.x[rs1]+imm, 1, cpu.x[rs2])

	case 0b_001: // sh
		success = cpu.memWrite(cpu.x[rs1]+imm, 2, cpu.x[rs2])

	case 0b_010: // sw
		success = cpu.memWrite(cpu.x[rs1]+imm, 4, cpu.x[rs2])

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if !success {
		cpu.trap(ExceptionStoreAMOAccessFault)
	}
}

func (cpu *CPU) execInstrFence(imm, rs1, f3, rd int32) {
	illegal := false

	switch f3 {
	case 0b_000: // fence
		if (imm&^0b_1111_1111) != 0 || rs1 != 0 || rd != 0 {
			illegal = true
		}

	case 0b_001: // fence.i
		if imm != 0 || rs1 != 0 || rd != 0 {
			illegal = true
		}

	default:
		illegal = true
	}

	if illegal {
		cpu.trap(ExceptionIllegalIstruction)
	}
}
