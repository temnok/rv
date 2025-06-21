package rv

const (
	ramBase = -1 << 31

	mepc = 0x341
)

type CPU struct {
	x   [32]int32
	pc  int32
	csr [4096]int32
	ram []byte

	instrIllegal bool
	accessFault  bool
	eCall        bool
}

func (cpu *CPU) executeInstr() {
	if cpu.faulted() {
		return
	}

	instr := cpu.readRAM(cpu.pc, 4)
	if cpu.faulted() {
		return
	}

	if instr&0b11 != 0b11 {
		cpu.instrIllegal = true
		return
	}

	rdBits := bits(instr, 7, 5)
	rd := &cpu.x[rdBits]
	f3 := bits(instr, 12, 3)
	rs1Bits := bits(instr, 15, 5)
	rs1 := cpu.x[rs1Bits]
	rs2 := cpu.x[bits(instr, 20, 5)]
	f7 := bits(instr, 25, 7)

	nextPC := cpu.pc + 4

	switch bits(instr, 2, 5) {
	case 0b_00000: // I-type loads
		switch f3 {
		case 0b_000: // lb
			if val := cpu.readRAM(rs1+immI(instr), 1); !cpu.faulted() {
				*rd = int32(int8(val))
			}
		case 0b_001: // lh
			if val := cpu.readRAM(rs1+immI(instr), 2); !cpu.faulted() {
				*rd = int32(int16(val))
			}
		case 0b_010: // lw
			if val := cpu.readRAM(rs1+immI(instr), 4); !cpu.faulted() {
				*rd = val
			}
		case 0b_100: // lbu
			if val := cpu.readRAM(rs1+immI(instr), 1); !cpu.faulted() {
				*rd = val
			}
		case 0b_101: // lhu
			if val := cpu.readRAM(rs1+immI(instr), 2); !cpu.faulted() {
				*rd = val
			}
		default:
			cpu.instrIllegal = true
		}

	case 0b_00011:
		switch f3 {
		case 0b_000: // fence
			if bits(instr, 28, 4) != 0 || rs1Bits != 0 || rdBits != 0 {
				cpu.instrIllegal = true
			}
		case 0b_001: // fence.i
			if bits(instr, 20, 12) != 0 || rs1Bits != 0 || rdBits != 0 {
				cpu.instrIllegal = true
			}
		}

	case 0b_00100: // I-type computational
		switch f3 {
		case 0b_000: // addi
			*rd = rs1 + immI(instr)
		case 0b_010: // slti
			if rs1 < immI(instr) {
				*rd = 1
			} else {
				*rd = 0
			}
		case 0b_011: // sltiu
			if uint32(rs1) < uint32(immI(instr)) {
				*rd = 1
			} else {
				*rd = 0
			}
		case 0b_100: // xori
			*rd = rs1 ^ immI(instr)
		case 0b_110: // ori
			*rd = rs1 | immI(instr)
		case 0b_111: // andi
			*rd = rs1 & immI(instr)
		default:
			switch f7<<3 | f3 {
			case 0b_0000000_001: // slli
				*rd = rs1 << immI(instr)
			case 0b_0000000_101: // srli
				*rd = int32(uint32(rs1) >> uint32(immI(instr)))
			case 0b_0100000_101: // srai
				*rd = rs1 >> (immI(instr) & 31)
			default:
				cpu.instrIllegal = true
			}
		}

	case 0b_00101: // U-type auipc
		*rd = cpu.pc + immU(instr)

	case 0b_01000: // S-type stores
		switch f3 {
		case 0b_000: // sb
			cpu.writeRAM(rs1+immS(instr), 1, rs2)
		case 0b_001: // sh
			cpu.writeRAM(rs1+immS(instr), 2, rs2)
		case 0b_010: // sw
			cpu.writeRAM(rs1+immS(instr), 4, rs2)
		default:
			cpu.instrIllegal = true
		}

	case 0b_01100: // R-type computational
		switch f7<<3 | f3 {
		case 0b_0000000_000: // add
			*rd = rs1 + rs2
		case 0b_0100000_000: // sub
			*rd = rs1 - rs2
		case 0b_0000000_001: // sll
			*rd = rs1 << (rs2 & 31)
		case 0b_0000000_010: // slt
			if rs1 < rs2 {
				*rd = 1
			} else {
				*rd = 0
			}
		case 0b_0000000_011: // sltu
			if uint32(rs1) < uint32(rs2) {
				*rd = 1
			} else {
				*rd = 0
			}
		case 0b_0000000_100: // xor
			*rd = rs1 ^ rs2
		case 0b_0000000_101: // srl
			*rd = int32(uint32(rs1) >> uint32(rs2&31))
		case 0b_0100000_101: // sra
			*rd = rs1 >> (rs2 & 31)
		case 0b_0000000_110: // or
			*rd = rs1 | rs2
		case 0b_0000000_111: // and
			*rd = rs1 & rs2
		default:
			cpu.instrIllegal = true
		}

	case 0b_01101: // U-type lui
		*rd = immU(instr)

	case 0b_11000: // B-type branches
		cond := false

		switch f3 {
		case 0b_000: // beq
			cond = rs1 == rs2
		case 0b_001: // bne
			cond = rs1 != rs2
		case 0b_100: // blt
			cond = rs1 < rs2
		case 0b_101: // bge
			cond = rs1 >= rs2
		case 0b_110: // bltu
			cond = uint32(rs1) < uint32(rs2)
		case 0b_111: // bgeu
			cond = uint32(rs1) >= uint32(rs2)
		default:
			cpu.instrIllegal = true
		}

		if cond {
			nextPC = cpu.pc + immB(instr)
		}

	case 0b_11001: // jalr
		*rd = nextPC
		nextPC = (rs1 + immI(instr)) &^ 1

	case 0b_11011: // jal
		*rd = nextPC
		nextPC = cpu.pc + immJ(instr)

	case 0b_11100: // I-type system
		imm := bits(instr, 20, 12)

		switch f3 {
		case 0b_000:
			if rdBits != 0 {
				cpu.instrIllegal = true
				break
			}
			switch imm {
			case 0b_0000_000_00000: // ecall
				cpu.eCall = true
			case 0b_0011_000_00010: // mret
				nextPC = cpu.csr[mepc]
			default:
				cpu.instrIllegal = true
			}
		case 0b_001: // csrrw
			*rd = cpu.csr[imm]
			cpu.csr[imm] = rs1
		case 0b_010: // csrrs
			*rd = cpu.csr[imm]
			cpu.csr[imm] |= rs1
		case 0b_011: // csrrc
			*rd = cpu.csr[imm]
			cpu.csr[imm] &^= rs1
		case 0b_101: // csrrwi
			*rd = cpu.csr[imm]
			cpu.csr[imm] = rs1Bits
		case 0b_110: // csrrsi
			*rd = cpu.csr[imm]
			cpu.csr[imm] |= rs1Bits
		case 0b_111: // csrrci
			*rd = cpu.csr[imm]
			cpu.csr[imm] &^= rs1Bits
		default:
			cpu.instrIllegal = true
		}

	default:
		cpu.instrIllegal = true
	}

	cpu.x[0] = 0

	if !cpu.faulted() {
		cpu.pc = nextPC
	}
}

func (cpu *CPU) faulted() bool {
	return cpu.instrIllegal || cpu.accessFault || cpu.eCall
}
