package rv

func (cpu *CPU) execInstrSystem(imm, rs1, f3, rd int32) {
	imm = bits(imm, 0, 12)

	switch f3 {
	case 0b_000:
		cpu.execInstrSystemSpecial(imm, rd)

	case 0b_001: // csrrw
		saved := cpu.csr[imm]
		cpu.csr[imm] = cpu.x[rs1]
		cpu.x[rd] = saved

	case 0b_010: // csrrs
		saved := cpu.csr[imm]
		cpu.csr[imm] |= cpu.x[rs1]
		cpu.x[rd] = saved

	case 0b_011: // csrrc
		saved := cpu.csr[imm]
		cpu.csr[imm] &^= cpu.x[rs1]
		cpu.x[rd] = saved

	case 0b_101: // csrrwi
		cpu.x[rd] = cpu.csr[imm]
		cpu.csr[imm] = rs1

	case 0b_110: // csrrsi
		cpu.x[rd] = cpu.csr[imm]
		cpu.csr[imm] |= rs1

	case 0b_111: // csrrci
		cpu.x[rd] = cpu.csr[imm]
		cpu.csr[imm] &^= rs1

	default:
		cpu.instrIllegal = true
	}
}

func (cpu *CPU) execInstrSystemSpecial(imm, rd int32) {
	if rd != 0 {
		cpu.instrIllegal = true
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		cpu.eCall = true

	case 0b_0011_000_00010: // mret
		cpu.nextPC = cpu.csr[mepc]

	default:
		cpu.instrIllegal = true
	}
}
