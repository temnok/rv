package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) fpDisabled() bool {
	return bi.Ts(cpu.CSR.Mstatus, csr.MstatusFS, 2) == csr.FSoff
}

func (cpu *CPU) execLoadFP(op instr.Op) {
	if cpu.fpDisabled() {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		return
	}

	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	switch op.F3() {
	case 0b_010: // flw
		if mem.Read(cpu.State, cpu.X[rs1]+imm, &val, 4); !trap.IsEntered(cpu.State) {
			cpu.Update.FVal = f32boxingBits | val
		}

	case 0b_011: // fld
		if !cpu.extD() {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
			return
		}

		if mem.Read(cpu.State, cpu.X[rs1]+imm, &val, 8); !trap.IsEntered(cpu.State) {
			cpu.Update.FVal = val
		}

	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		return
	}

	cpu.Update.FReg = rd
}

func (cpu *CPU) execStoreFP(op instr.Op) {
	if cpu.fpDisabled() {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		return
	}

	imm, rs1, rs2 := imm.S(op.Code()), op.Rs1(), op.Rs2()

	switch op.F3() {
	case 0b_010: // fsw
		mem.Write(cpu.State, cpu.X[rs1]+imm, cpu.F[rs2], 4)

	case 0b_011: // fsd
		if !cpu.extD() {
			trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
			return
		}

		mem.Write(cpu.State, cpu.X[rs1]+imm, cpu.F[rs2], 8)

	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}
