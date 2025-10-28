package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func (cpu *CPU) fpDisabled() bool {
	return bi.Ts(cpu.CSR.Mstatus, state.MstatusFS, 2) == state.FSoff
}

func (cpu *CPU) execLoadFP(imm, rs1, f3, rd int) {
	if cpu.fpDisabled() {
		cpu.Trap(ExceptionIllegalIstruction)
		return
	}

	var val int

	switch f3 {
	case 0b_010: // flw
		if cpu.memRead(cpu.X[rs1]+imm, &val, 4); !cpu.IsTrapped() {
			cpu.Update.FVal = f32boxingBits | val
		}

	case 0b_011: // fld
		if !cpu.extD() {
			cpu.Trap(ExceptionIllegalIstruction)
			return
		}

		if cpu.memRead(cpu.X[rs1]+imm, &val, 8); !cpu.IsTrapped() {
			cpu.Update.FVal = val
		}

	default:
		cpu.Trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Update.FReg = rd
}

func (cpu *CPU) execStoreFP(imm, rs2, rs1, f3 int) {
	if cpu.fpDisabled() {
		cpu.Trap(ExceptionIllegalIstruction)
		return
	}

	switch f3 {
	case 0b_010: // fsw
		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 4)

	case 0b_011: // fsd
		if !cpu.extD() {
			cpu.Trap(ExceptionIllegalIstruction)
			return
		}

		cpu.memWrite(cpu.X[rs1]+imm, cpu.F[rs2], 8)

	default:
		cpu.Trap(ExceptionIllegalIstruction)
	}
}
