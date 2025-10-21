package rv

import "github.com/temnok/rv/instr"

func (cpu *CPU) execComputeI(imm, rs1, f3, rd int) {
	switch f3 {
	case 0b_000:
		instr.Addi(&cpu.State, rd, rs1, imm)
	case 0b_001:
		instr.Slli(&cpu.State, rd, rs1, imm)
	case 0b_010:
		instr.Slti(&cpu.State, rd, rs1, imm)

	case 0b_011: // sltiu
		if cpu.Xuint(cpu.X[rs1]) < cpu.Xuint(imm) {
			cpu.XRegSet(rd, 1)
		} else {
			cpu.XRegSet(rd, 0)
		}

	case 0b_100: // xori
		cpu.XRegSet(rd, cpu.X[rs1]^imm)

	case 0b_101:
		if imm < cpu.XLen { // srli
			cpu.XRegSet(rd, int(cpu.Xuint(cpu.X[rs1])>>cpu.Xuint(imm)))
		} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
			cpu.XRegSet(rd, cpu.X[rs1]>>imm)
		}

	case 0b_110: // ori
		cpu.XRegSet(rd, cpu.X[rs1]|imm)

	case 0b_111: // andi
		cpu.XRegSet(rd, cpu.X[rs1]&imm)
	}

	if cpu.Update.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}

func (cpu *CPU) execComputeR(f7, rs2, rs1, f3, rd int) {
	if f7 == 1 {
		cpu.execComputeM(rs2, rs1, f3, rd)
		return
	}

	op := bit(f7, 5)<<3 | f3
	if f7 &^= 0b0100000; f7 != 0 {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	switch op {
	case 0b_000: // add
		cpu.XRegSet(rd, cpu.X[rs1]+cpu.X[rs2])

	case 0b_1_000: // sub
		cpu.XRegSet(rd, cpu.X[rs1]-cpu.X[rs2])

	case 0b_001: // sll
		cpu.XRegSet(rd, cpu.X[rs1]<<(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_010: // slt
		if cpu.X[rs1] < cpu.X[rs2] {
			cpu.XRegSet(rd, 1)
		} else {
			cpu.XRegSet(rd, 0)
		}

	case 0b_011: // sltu
		if cpu.Xuint(cpu.X[rs1]) < cpu.Xuint(cpu.X[rs2]) {
			cpu.XRegSet(rd, 1)
		} else {
			cpu.XRegSet(rd, 0)
		}

	case 0b_100: // xor
		cpu.XRegSet(rd, cpu.X[rs1]^cpu.X[rs2])

	case 0b_101: // srl
		cpu.XRegSet(rd, int(cpu.Xuint(cpu.X[rs1])>>cpu.Xuint(cpu.X[rs2]&(cpu.XLen-1))))

	case 0b_1_101: // sra
		cpu.XRegSet(rd, cpu.X[rs1]>>(cpu.X[rs2]&(cpu.XLen-1)))

	case 0b_110: // or
		cpu.XRegSet(rd, cpu.X[rs1]|cpu.X[rs2])

	case 0b_111: // and
		cpu.XRegSet(rd, cpu.X[rs1]&cpu.X[rs2])
	}

	if cpu.Update.XReg < 0 {
		cpu.trap(ExceptionIllegalIstruction)
	}
}
