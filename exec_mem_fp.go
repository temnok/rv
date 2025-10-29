package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) fpDisabled() bool {
	return bi.Ts(cpu.CSR.Mstatus, csr.MstatusFS, 2) == csr.FSoff
}

func (cpu *CPU) execLoadFP(imm, rs1, f3, rd int) {
	if cpu.fpDisabled() {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	var val int

	switch f3 {
	case 0b_010: // flw
		if cpu.memRead(cpu.X[rs1]+imm, &val, 4); !trap.IsEntered(cpu.State) {
			cpu.Update.FVal = f32boxingBits | val
		}

	case 0b_011: // fld
		if !cpu.extD() {
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.X[rs1]+imm, &val, 8); !trap.IsEntered(cpu.State) {
			cpu.Update.FVal = val
		}

	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	cpu.Update.FReg = rd
}

func (cpu *CPU) execStoreFP(imm, rs2, rs1, f3 int) {
	if cpu.fpDisabled() {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	switch f3 {
	case 0b_010: // fsw
		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 4)

	case 0b_011: // fsd
		if !cpu.extD() {
			trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 8)

	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
