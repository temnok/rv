package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execSystem(imm, rs1, f3, rd int) {
	if f3 == 0 {
		cpu.execSystemSpecial(imm, rd)
	} else {
		cpu.execSystemCSR(imm, rs1, f3, rd)
	}
}

func (cpu *CPU) execSystemSpecial(imm, rd int) {
	if rd != 0 {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		trap.EnterWithoutTval(cpu.State, ExceptionEnvironmentCallFromUMode+cpu.Priv)

	case 0b_0000_000_00001: // ebreak
		trap.EnterWithoutTval(cpu.State, ExceptionBreakpoint)

	case 0b_0001_000_00010: // sret
		trap.Exit(cpu.State, PrivS)

	case 0b_0001_000_00101: // wfi, https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#wfi

	case 0b_0011_000_00010: // mret
		trap.Exit(cpu.State, PrivM)

	default:
		switch bi.Ts(imm, 5, 7) {
		case 0b_0001_001: // sfence.vma
			cpu.TLB.flush()
			cpu.Update.ICache.Clear()

			if cpu.Priv == PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTVM) == 1 {
				trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			}

		default:
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		}
	}
}

func (cpu *CPU) execSystemCSR(imm, rs1, f3, rd int) {
	csrReg := bi.Ts(imm, 0, 12)

	s := rs1
	if (f3 & 0b_100) == 0 {
		s = cpu.X[s]
	}

	var val int

	switch f3 & 3 {
	case 0b_01: // csrrw
		if rd != 0 {
			if !csr.Read(cpu.State, csrReg, &val) {
				trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
				return
			}
		}

		if !csr.Write(cpu.State, csrReg, s) {
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			return
		}
		cpu.Xset(rd, val)

	case 0b_10: // csrrs
		if !csr.Read(cpu.State, csrReg, &val) {
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			return
		}

		if s != 0 {
			if !csr.Write(cpu.State, csrReg, val|s) {
				trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
				return
			}
		}

		cpu.Xset(rd, val)

	case 0b_11: // csrrc
		if !csr.Read(cpu.State, csrReg, &val) {
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			return
		}

		if s != 0 {
			if !csr.Write(cpu.State, csrReg, val&^s) {
				trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
				return
			}
		}

		cpu.Xset(rd, val)

	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
