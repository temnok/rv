package rv

func (cpu *CPU) exec(opcode int) {
	f7 := bits(opcode, 25, 7)
	rs2 := bits(opcode, 20, 5)
	rs1 := bits(opcode, 15, 5)
	f3 := bits(opcode, 12, 3)
	rd := bits(opcode, 7, 5)

	switch bits(opcode, 2, 5) {
	case 0b_00000:
		cpu.execLoad(immI(opcode), rs1, f3, rd)

	case 0b_00011:
		cpu.execFence(immI(opcode), rs1, f3, rd)

	case 0b_00100:
		cpu.execComputeI(immI(opcode), rs1, f3, rd)

	case 0b_00110:
		cpu.execComputeIw(immI(opcode), rs1, f3, rd)

	case 0b_00101: // auipc
		cpu.x[rd] = cpu.xint(cpu.pc + immU(opcode))

	case 0b_01000:
		cpu.execStore(immS(opcode), rs2, rs1, f3)

	case 0b_01011:
		cpu.execAtomic(f7, rs2, rs1, f3, rd)

	case 0b_01100:
		cpu.execComputeR(f7, rs2, rs1, f3, rd)

	case 0b_01110:
		cpu.execComputeRw(f7, rs2, rs1, f3, rd)

	case 0b_01101: // lui
		cpu.x[rd] = immU(opcode)

	case 0b_11000:
		cpu.execBranch(immB(opcode), rs2, rs1, f3)

	case 0b_11001: // jalr
		saved := cpu.next.pc
		cpu.next.pc = cpu.xint((cpu.x[rs1] + immI(opcode)) &^ 1)
		cpu.x[rd] = saved

	case 0b_11011: // jal
		cpu.x[rd] = cpu.next.pc
		cpu.next.pc = cpu.xint(cpu.pc + immJ(opcode))

	case 0b_11100:
		cpu.execSystem(immI(opcode), rs1, f3, rd)

	default:
		cpu.trap(ExceptionIllegalIstruction)
	}

	cpu.x[0] = 0
}
