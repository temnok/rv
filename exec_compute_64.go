package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/trap"
)

func (cpu *CPU) execComputeI64(op instr.Op) {
	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	if cpu.Xlen64() {
		switch op.F3() {
		case 0b_000:
			instr.Addiw(cpu.State, rd, rs1, imm)
		case 0b_001:
			if imm < 32 {
				instr.Slliw(cpu.State, rd, rs1, imm)
			}
		case 0b_101:
			if imm < 32 {
				instr.Srliw(cpu.State, rd, rs1, imm)
			} else if imm &^= 0b0100000_00000; imm < 32 {
				instr.Sraiw(cpu.State, rd, rs1, imm)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execComputeR64(op instr.Op) {
	if !cpu.Xlen64() {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	f7, f3, rd, rs1, rs2 := op.F7(), op.F3(), op.Rd(), op.Rs1(), op.Rs2()

	if f7 == 1 {
		cpu.execComputeM64(rs2, rs1, f3, rd)
		return
	}

	if f7&^0b0100000 != 0 {
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
		return
	}

	switch bi.T(f7, 5)<<3 | f3 {
	case 0b_000:
		instr.Addw(cpu.State, rd, rs1, rs2)
	case 0b_1_000:
		instr.Subw(cpu.State, rd, rs1, rs2)
	case 0b_001:
		instr.Sllw(cpu.State, rd, rs1, rs2)
	case 0b_101:
		instr.Srlw(cpu.State, rd, rs1, rs2)
	case 0b_1_101:
		instr.Sraw(cpu.State, rd, rs1, rs2)
	default:
		trap.EnterWithoutTval(cpu.State, ExceptionIllegalIstruction)
	}
}
