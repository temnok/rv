package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execComputeI64(op instr.Op) {
	imm := imm.I(op.Code())

	if cpu.Xlen64() {
		switch op.F3() {
		case 0b_000:
			instr.Addiw(cpu.State, op)
		case 0b_001:
			if imm < 32 {
				instr.Slliw(cpu.State, op)
			}
		case 0b_101:
			if imm < 32 {
				instr.Srliw(cpu.State, op)
			} else if imm&^0b0100000_00000 < 32 {
				instr.Sraiw(cpu.State, op)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}

func (cpu *CPU) execComputeR64(op instr.Op) {
	if !cpu.Xlen64() {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		return
	}

	f7, f3 := op.F7(), op.F3()

	if f7 == 1 {
		cpu.execComputeM64(op)
		return
	}

	if f7&^0b0100000 != 0 {
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
		return
	}

	switch bi.T(f7, 5)<<3 | f3 {
	case 0b_000:
		instr.Addw(cpu.State, op)
	case 0b_1_000:
		instr.Subw(cpu.State, op)
	case 0b_001:
		instr.Sllw(cpu.State, op)
	case 0b_101:
		instr.Srlw(cpu.State, op)
	case 0b_1_101:
		instr.Sraw(cpu.State, op)
	default:
		trap.EnterWithoutTval(cpu.State, trap.IllegalIstruction)
	}
}
