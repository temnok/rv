package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/instr"
)

func (cpu *CPU) execComputeI64(imm, rs1, f3, rd int) {
	if cpu.Xlen64() {
		switch f3 {
		case 0b_000:
			instr.Addiw(&cpu.State, rd, rs1, imm)
		case 0b_001:
			if imm < 32 {
				instr.Slliw(&cpu.State, rd, rs1, imm)
			}
		case 0b_101:
			if imm < 32 {
				instr.Srliw(&cpu.State, rd, rs1, imm)
			} else if imm &^= 0b0100000_00000; imm < 32 {
				instr.Sraiw(&cpu.State, rd, rs1, imm)
			}
		}
	}

	if cpu.Update.XReg < 0 {
		cpu.Trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execComputeR64(f7, rs2, rs1, f3, rd int) {
	if !cpu.Xlen64() {
		cpu.Trap(ExceptionIllegalIstruction)
		return
	}

	if f7 == 1 {
		cpu.execComputeM64(rs2, rs1, f3, rd)
		return
	}

	op := bi.T(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.Trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000:
		instr.Addw(&cpu.State, rd, rs1, rs2)
	case 0b_1_000:
		instr.Subw(&cpu.State, rd, rs1, rs2)
	case 0b_001:
		instr.Sllw(&cpu.State, rd, rs1, rs2)
	case 0b_101:
		instr.Srlw(&cpu.State, rd, rs1, rs2)
	case 0b_1_101:
		instr.Sraw(&cpu.State, rd, rs1, rs2)
	default:
		cpu.Trap(ExceptionIllegalIstruction)
	}
}
