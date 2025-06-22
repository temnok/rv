package rv

func (cpu *CPU) execInstrSystem(imm, rs1, f3, rd int32) {
	imm = bits(imm, 0, 12)

	if f3 == 0 {
		cpu.execInstrSystemSpecial(imm, rd)
	} else {
		cpu.execInstrSystemCSR(imm, rs1, f3, rd)
	}
}

func (cpu *CPU) execInstrSystemSpecial(imm, rd int32) {
	if rd != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		cpu.trap(ExceptionEnvironmentCallFromMMode)

	case 0b_0011_000_00010: // mret
		cpu.nextPC = cpu.csr.mepc

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execInstrSystemCSR(imm, rs1, f3, rd int32) {
	s := rs1
	if (f3 & 0b_100) == 0 {
		s = cpu.x[s]
	}

	var val int32

	switch f3 & 0b_11 {
	case 0b_01: // csrrw
		if rd == 0 || cpu.csrRead(imm, &val) {
			if cpu.csrWrite(imm, s) {
				cpu.x[rd] = val
			}
		}

	case 0b_10: // csrrs
		if cpu.csrRead(imm, &val) {
			if s == 0 || cpu.csrWrite(imm, val|s) {
				cpu.x[rd] = val
			}
		}

	case 0b_11: // csrrc
		if cpu.csrRead(imm, &val) {
			if s == 0 || cpu.csrWrite(imm, val&^s) {
				cpu.x[rd] = val
			}
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
