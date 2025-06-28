package rv

func (cpu *CPU) execSystem(imm, rs1, f3, rd Xint) {
	if f3 == 0 {
		cpu.execSystemSpecial(imm, rd)
	} else {
		cpu.execSystemCSR(imm, rs1, f3, rd)
	}
}

func (cpu *CPU) execSystemSpecial(imm, rd Xint) {
	if rd != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		cpu.trap(ExceptionEnvironmentCallFromUMode + cpu.priv)

	case 0b_0000_000_00001: // ebreak
		cpu.trap(ExceptionBreakpoint)

	case 0b_0001_000_00010: // sret
		cpu.ret(PrivS)

	case 0b_0001_000_00101: // wfi, https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#wfi

	case 0b_0011_000_00010: // mret
		cpu.ret(PrivM)

	default:
		switch bits(imm, 5, 7) {
		case 0b_0001_001: // sfence.vma
			cpu.tlb.flush()

			if cpu.priv == PrivS && bit(cpu.csr.mstatus, mstatusTVM) != 0 {
				cpu.trap(ExceptionIllegalIstruction)
			}

		default:
			cpu.trap(ExceptionIllegalIstruction)
		}
	}
}

func (cpu *CPU) execSystemCSR(imm, rs1, f3, rd Xint) {
	csr := bits(imm, 0, 12)

	s := rs1
	if (f3 & 0b_100) == 0 {
		s = cpu.x[s]
	}

	var val Xint

	switch f3 & 0b_11 {
	case 0b_01: // csrrw
		if (rd == 0 || cpu.csrRead(csr, &val)) && cpu.csrWrite(csr, s) {
			cpu.x[rd] = val
		}

	case 0b_10: // csrrs
		if cpu.csrRead(csr, &val) && (s == 0 || cpu.csrWrite(csr, val|s)) {
			cpu.x[rd] = val
		}

	case 0b_11: // csrrc
		if cpu.csrRead(csr, &val) && (s == 0 || cpu.csrWrite(csr, val&^s)) {
			cpu.x[rd] = val
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
