package rv

func (cpu *CPU) execLoadFP(imm, rs1, f3, rd int) {
	if bits(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	var val int

	switch f3 {
	case 0b_010: // flw
		if cpu.memRead(cpu.X[rs1]+imm, &val, 4); !cpu.isTrapped() {
			cpu.Updated.FVal = f32boxingBits | val
		}

	case 0b_011: // fld
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.X[rs1]+imm, &val, 8); !cpu.isTrapped() {
			cpu.Updated.FVal = val
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.FReg = rd
}

func (cpu *CPU) execStoreFP(imm, rs2, rs1, f3 int) {
	if bits(cpu.CSR.Mstatus, MstatusFS, 2) == FSoff {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch f3 {
	case 0b_010: // fsw
		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 4)

	case 0b_011: // fsd
		if !cpu.xlen64() {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 8)

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
