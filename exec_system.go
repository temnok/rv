package rv

func (cpu *CPU) execSystem(imm, rs1, f3, rd int) {
	if f3 == 0 {
		cpu.execSystemSpecial(imm, rd)
	} else {
		cpu.execSystemCSR(imm, rs1, f3, rd)
	}
}

func (cpu *CPU) execSystemSpecial(imm, rd int) {
	if rd != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		cpu.trap(ExceptionEnvironmentCallFromUMode + cpu.Priv)

	case 0b_0000_000_00001: // ebreak
		cpu.trap(ExceptionBreakpoint)

	case 0b_0001_000_00010: // sret
		cpu.xret(PrivS)

	case 0b_0001_000_00101: // wfi, https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#wfi

	case 0b_0011_000_00010: // mret
		cpu.xret(PrivM)

	default:
		switch bits(imm, 5, 7) {
		case 0b_0001_001: // sfence.vma
			cpu.TLB.flush()
			cpu.Updated.ICache.Clear()

			if cpu.Priv == PrivS && bit(cpu.CSR.Mstatus, MstatusTVM) == 1 {
				cpu.trap(ExceptionIllegalIstruction)
			}

		default:
			cpu.trap(ExceptionIllegalIstruction)
		}
	}
}

func (cpu *CPU) execSystemCSR(imm, rs1, f3, rd int) {
	csr := bits(imm, 0, 12)

	s := rs1
	if (f3 & 0b_100) == 0 {
		s = cpu.Reg[s]
	}

	var val int

	switch f3 & 3 {
	case 0b_01: // csrrw
		if rd != 0 {
			if cpu.csrRead(csr, &val); cpu.isTrapped() {
				return
			}
		}

		if cpu.csrWrite(csr, s); !cpu.isTrapped() {
			cpu.Updated.RegIndex = rd
			cpu.Updated.RegValue = val
		}

	case 0b_10: // csrrs
		if cpu.csrRead(csr, &val); cpu.isTrapped() {
			return
		}

		if s != 0 {
			if cpu.csrWrite(csr, val|s); cpu.isTrapped() {
				return
			}
		}

		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = val

	case 0b_11: // csrrc
		if cpu.csrRead(csr, &val); cpu.isTrapped() {
			return
		}

		if s != 0 {
			if cpu.csrWrite(csr, val&^s); cpu.isTrapped() {
				return
			}
		}

		cpu.Updated.RegIndex = rd
		cpu.Updated.RegValue = val

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
