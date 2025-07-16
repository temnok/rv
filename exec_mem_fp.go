package rv

func (cpu *CPU) execLoadFP(imm, rs1, f3, rd int) {
	var val int

	switch f3 {
	case 0b_010: // flw
		if cpu.memRead(cpu.Reg[rs1]+imm, &val, 4); !cpu.isTrapped() {
			cpu.Updated.FRegValue = upper32ones | val
		}

	case 0b_011: // fld
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.Reg[rs1]+imm, &val, 8); !cpu.isTrapped() {
			cpu.Updated.FRegValue = val
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.FRegUpdated = true
	cpu.Updated.FRegIndex = rd
}

func (cpu *CPU) execStoreFP(imm, rs2, rs1, f3 int) {
	switch f3 {
	case 0b_010: // fsw
		cpu.memWrite(cpu.Reg[rs1]+imm, cpu.FReg[rs2], 4)

	case 0b_011: // fsd
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.Reg[rs1]+imm, cpu.FReg[rs2], 8)

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
