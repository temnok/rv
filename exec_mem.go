package rv

func (cpu *CPU) execLoad(imm, rs1, f3, rd int32) {
	var val int32

	switch f3 {
	case 0b_000: // lb
		if cpu.memRead(cpu.x[rs1]+imm, 1, &val) {
			cpu.x[rd] = int32(int8(val))
		}

	case 0b_001: // lh
		if cpu.memRead(cpu.x[rs1]+imm, 2, &val) {
			cpu.x[rd] = int32(int16(val))
		}

	case 0b_010: // lw
		if cpu.memRead(cpu.x[rs1]+imm, 4, &val) {
			cpu.x[rd] = val
		}

	case 0b_100: // lbu
		if cpu.memRead(cpu.x[rs1]+imm, 1, &val) {
			cpu.x[rd] = int32(uint8(val))
		}

	case 0b_101: // lhu
		if cpu.memRead(cpu.x[rs1]+imm, 2, &val) {
			cpu.x[rd] = int32(uint16(val))
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execStore(imm, rs2, rs1, f3 int32) {
	switch f3 {
	case 0b_000: // sb
		cpu.memWrite(cpu.x[rs1]+imm, 1, int32(uint8(cpu.x[rs2])))

	case 0b_001: // sh
		cpu.memWrite(cpu.x[rs1]+imm, 2, int32(uint16(cpu.x[rs2])))

	case 0b_010: // sw
		cpu.memWrite(cpu.x[rs1]+imm, 4, cpu.x[rs2])

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execFence(imm, rs1, f3, rd int32) {
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
