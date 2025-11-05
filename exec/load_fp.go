package exec

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func LoadFP(cpu *state.State, op instr.Op) {
	if csr.FpDisabled(cpu) {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	switch op.F3() {
	case 0b_010: // flw
		if mem.Read(cpu, cpu.X[rs1]+imm, &val, 4); !trap.IsEntered(cpu) {
			cpu.Update.FVal = f32boxingBits | val
		}

	case 0b_011: // fld
		if !csr.ExtD(cpu) {
			trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
			return
		}

		if mem.Read(cpu, cpu.X[rs1]+imm, &val, 8); !trap.IsEntered(cpu) {
			cpu.Update.FVal = val
		}

	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	cpu.Update.FReg = rd
}
