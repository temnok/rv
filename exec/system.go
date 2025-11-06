package exec

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func System(cpu *state.CPU, op instr.Op) {
	if op.F3() == 0 {
		systemSpecial(cpu, op)
	} else {
		systemCSR(cpu, op)
	}
}

func systemSpecial(cpu *state.CPU, op instr.Op) {
	imm, rd := imm.I(op.Code()), op.Rd()

	if rd != 0 {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	switch imm {
	case 0b_0000_000_00000: // ecall
		trap.EnterWithoutTval(cpu, trap.EnvironmentCallFromUMode+cpu.Priv)

	case 0b_0000_000_00001: // ebreak
		trap.EnterWithoutTval(cpu, trap.Breakpoint)

	case 0b_0001_000_00010: // sret
		trap.Exit(cpu, state.PrivS)

	case 0b_0001_000_00101: // wfi, https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#wfi

	case 0b_0011_000_00010: // mret
		trap.Exit(cpu, state.PrivM)

	default:
		switch bi.Ts(imm, 5, 7) {
		case 0b_0001_001: // sfence.vma
			cpu.TLB.Flush()
			cpu.Update.ICache.Clear()

			if cpu.Priv == state.PrivS && bi.T(cpu.CSR.Mstatus, csr.MstatusTVM) == 1 {
				trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			}

		default:
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		}
	}
}

func systemCSR(cpu *state.CPU, op instr.Op) {
	switch op.F3() & 3 {
	case 1:
		instr.Csrrw(cpu, op)
	case 2:
		instr.Csrrs(cpu, op)
	case 3:
		instr.Csrrc(cpu, op)
	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
	}
}
