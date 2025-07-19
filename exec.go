package rv

func (cpu *CPU) exec(opcode int) {
	opcodeSize := 4
	if isCompressed := opcode&3 != 3; isCompressed {
		opcodeSize = 2

		if cpu.decompress(&opcode); cpu.isTrapped() {
			return
		}
	}
	cpu.Updated.PC = cpu.xint(cpu.PC + opcodeSize)

	f7 := bits(opcode, 25, 7)
	rs2 := bits(opcode, 20, 5)
	rs1 := bits(opcode, 15, 5)
	f3 := bits(opcode, 12, 3)
	rd := bits(opcode, 7, 5)

	switch op := bits(opcode, 2, 5); op {
	case 0b_00000:
		cpu.execLoad(immI(opcode), rs1, f3, rd)

	case 0b_00001:
		cpu.execLoadFP(immI(opcode), rs1, f3, rd)

	case 0b_00011:
		cpu.execFence(immI(opcode), rs1, f3, rd)

	case 0b_00100:
		cpu.execComputeI(immI(opcode), rs1, f3, rd)

	case 0b_00110:
		cpu.execComputeI64(immI(opcode), rs1, f3, rd)

	case 0b_00101: // auipc
		cpu.Updated.XReg = rd
		cpu.Updated.XVal = cpu.xint(cpu.PC + immU(opcode))

	case 0b_01000:
		cpu.execStore(immS(opcode), rs2, rs1, f3)

	case 0b_01001:
		cpu.execStoreFP(immS(opcode), rs2, rs1, f3)

	case 0b_01011:
		cpu.execAtomic(f7, rs2, rs1, f3, rd)

	case 0b_01100:
		cpu.execComputeR(f7, rs2, rs1, f3, rd)

	case 0b_01110:
		cpu.execComputeR64(f7, rs2, rs1, f3, rd)

	case 0b_01101: // lui
		cpu.Updated.XReg = rd
		cpu.Updated.XVal = immU(opcode)

	case 0b_10000, 0b_10001, 0b_10010, 0b_10011:
		cpu.execComputeFP(f7, rs2, rs1, f3, rd, op)

	case 0b_10100:
		cpu.execComputeFP(f7, rs2, rs1, f3, rd, 0)

	case 0b_11000:
		cpu.execBranch(immB(opcode), rs2, rs1, f3)

	case 0b_11001: // jalr
		cpu.Updated.XReg = rd
		cpu.Updated.XVal = cpu.PC + opcodeSize
		cpu.Updated.PC = cpu.xint((cpu.X[rs1] + immI(opcode)) &^ 1)

	case 0b_11011: // jal
		cpu.Updated.XReg = rd
		cpu.Updated.XVal = cpu.PC + opcodeSize
		cpu.Updated.PC = cpu.xint(cpu.PC + immJ(opcode))

	case 0b_11100:
		cpu.execSystem(immI(opcode), rs1, f3, rd)

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}
}
