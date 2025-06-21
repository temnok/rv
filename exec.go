package rv

func (cpu *CPU) execInstr() {
	if cpu.faulted() {
		return
	}

	opcode := cpu.memRead(cpu.pc, 4)
	if cpu.faulted() {
		return
	}

	if opcode&0b11 == 0b11 {
		cpu.nextPC = cpu.pc + 4
	} else {
		cpu.nextPC = cpu.pc + 2

		opcode = decompress(opcode)
	}

	f7 := bits(opcode, 25, 7)
	rs2 := bits(opcode, 20, 5)
	rs1 := bits(opcode, 15, 5)
	f3 := bits(opcode, 12, 3)
	rd := bits(opcode, 7, 5)

	switch bits(opcode, 2, 5) {
	case 0b_00000:
		cpu.execInstrLoad(immI(opcode), rs1, f3, rd)

	case 0b_00011:
		cpu.execInstrFence(immI(opcode), rs1, f3, rd)

	case 0b_00100:
		cpu.execInstrComputeI(immI(opcode), rs1, f3, rd)

	case 0b_00101: // auipc
		cpu.x[rd] = cpu.pc + immU(opcode)

	case 0b_01000:
		cpu.execInstrStore(immS(opcode), rs2, rs1, f3)

	case 0b_01100:
		cpu.execInstrComputeR(f7, rs2, rs1, f3, rd)

	case 0b_01101: // lui
		cpu.x[rd] = immU(opcode)

	case 0b_11000:
		cpu.execInstrBranch(immB(opcode), rs2, rs1, f3)

	case 0b_11001: // jalr
		saved := cpu.nextPC
		cpu.nextPC = (cpu.x[rs1] + immI(opcode)) &^ 1
		cpu.x[rd] = saved

	case 0b_11011: // jal
		cpu.x[rd] = cpu.nextPC
		cpu.nextPC = cpu.pc + immJ(opcode)

	case 0b_11100:
		cpu.execInstrSystem(immI(opcode), rs1, f3, rd)

	default:
		cpu.instrIllegal = true
	}

	cpu.x[0] = 0

	if !cpu.faulted() {
		cpu.pc = cpu.nextPC
	}
}
