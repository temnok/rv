package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/instr"
)

func (cpu *CPU) exec(opcode int) {
	opcodeSize := 4
	if isCompressed := opcode&3 != 3; isCompressed {
		opcodeSize = 2

		if cpu.decompress(&opcode); cpu.IsTrapped() {
			return
		}
	}
	cpu.Update.PC = cpu.Xint(cpu.PC + opcodeSize)

	f7 := bi.Ts(opcode, 25, 7)
	rs2 := bi.Ts(opcode, 20, 5)
	rs1 := bi.Ts(opcode, 15, 5)
	f3 := bi.Ts(opcode, 12, 3)
	rd := bi.Ts(opcode, 7, 5)

	switch op := bi.Ts(opcode, 2, 5); op {
	case 0b_00000:
		cpu.execLoad(imm.I(opcode), rs1, f3, rd)
	case 0b_00001:
		cpu.execLoadFP(imm.I(opcode), rs1, f3, rd)
	case 0b_00011:
		cpu.execFence(imm.I(opcode), rs1, f3, rd)
	case 0b_00100:
		cpu.execComputeI(imm.I(opcode), rs1, f3, rd)
	case 0b_00110:
		cpu.execComputeI64(imm.I(opcode), rs1, f3, rd)
	case 0b_00101:
		instr.Auipc(&cpu.State, rd, imm.U(opcode))
	case 0b_01000:
		cpu.execStore(imm.S(opcode), rs2, rs1, f3)
	case 0b_01001:
		cpu.execStoreFP(imm.S(opcode), rs2, rs1, f3)
	case 0b_01011:
		cpu.execAtomic(f7, rs2, rs1, f3, rd)
	case 0b_01100:
		cpu.execComputeR(f7, rs2, rs1, f3, rd)
	case 0b_01110:
		cpu.execComputeR64(f7, rs2, rs1, f3, rd)
	case 0b_01101:
		instr.Lui(&cpu.State, rd, imm.U(opcode))
	case 0b_10000, 0b_10001, 0b_10010, 0b_10011:
		cpu.execComputeFP(f7, rs2, rs1, f3, rd, op)
	case 0b_10100:
		cpu.execComputeFP(f7, rs2, rs1, f3, rd, 0)
	case 0b_11000:
		cpu.execBranch(imm.B(opcode), rs2, rs1, f3)
	case 0b_11001:
		instr.Jalr(&cpu.State, opcodeSize, rd, rs1, imm.I(opcode))
	case 0b_11011:
		instr.Jal(&cpu.State, opcodeSize, rd, imm.J(opcode))
	case 0b_11100:
		cpu.execSystem(imm.I(opcode), rs1, f3, rd)
	default:
		cpu.Trap(ExceptionIllegalIstruction)
	}
}
