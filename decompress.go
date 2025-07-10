package rv

func (cpu *CPU) decompress(opcodePtr *int) {
	opcode := *opcodePtr
	if is32bitInstruction := opcode&3 == 3; is32bitInstruction {
		return
	}

	opcode = int(uint16(opcode))
	decompressedOpcode := decompress(opcode)
	if decompressedOpcode == 0 {
		cpu.trapWithTval(ExceptionIllegalIstruction, opcode)
		return
	}

	*opcodePtr = decompressedOpcode
	cpu.nextPC = cpu.pc + 2
}

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rvc_instruction_set_listings
func decompress(opcode int) int {
	f3 := bits(opcode, 13, 3)
	rA := bits(opcode, 7, 5)
	ra := 8 + (rA & 7)
	rB := bits(opcode, 2, 5)
	rb := 8 + (rB & 7)

	switch op := bits(opcode, 0, 2); op {
	case 0b_00:
		switch f3 {
		case 0b_000:
			if imm := immCIW(opcode); imm != 0 {
				return encodeI(imm, 2, 0, rb, 4) // addi
			}

		case 0b_010: // c.lw
			return encodeI(immCL(opcode), ra, 2, rb, 0) // lw

		case 0b_011: // c.ld
			if Xlen64 {
				return encodeI(immCLx8(opcode), ra, 3, rb, 0) // ld
			}

		case 0b_110: // c.sw
			return encodeS(immCL(opcode), rb, ra, 2, 8) // sw

		case 0b_111: // c.sd
			if Xlen64 {
				return encodeS(immCLx8(opcode), rb, ra, 3, 8) // sw
			}
		}

	case 0b_01:
		switch f3 {
		case 0b_000: // c.addi
			return encodeI(immCI(opcode), rA, 0, rA, 4)

		case 0b_001:
			if Xlen32 {
				return encodeJ(immCJ(opcode), 1, 27) // jal
			} else if rA != 0 {
				return encodeI(immCI(opcode), rA, 0, rA, 6) // addiw
			}

		case 0b_010: // li
			return encodeI(immCI(opcode), 0, 0, rA, 4) // addi

		case 0b_011:
			switch rA {
			case 0: // illegal
				break

			case 2: // c.addi16sp
				return encodeI(immCIx16(opcode), 2, 0, 2, 4)

			default: // c.lui
				return encodeU(immCI(opcode), rA, 13)
			}

		case 0b_100:
			switch bits(opcode, 10, 2) {
			case 0b_00: // srli
				return encodeR(0, immCI(opcode)&(Xlen-1), ra, 5, ra, 4)

			case 0b_01: // srai
				return encodeR(0b_0100000, immCI(opcode)&(Xlen-1), ra, 5, ra, 4)

			case 0b_10: // andi
				return encodeI(immCI(opcode), ra, 7, ra, 4)

			case 0b_11:
				switch bit(opcode, 12)<<2 | bits(opcode, 5, 2) {
				case 0b_000: // c.sub
					return encodeR(0b_0100000, rb, ra, 0, ra, 12)

				case 0b_001: // c.xor
					return encodeR(0, rb, ra, 4, ra, 12)

				case 0b_010: // c.or
					return encodeR(0, rb, ra, 6, ra, 12)

				case 0b_011: // c.and
					return encodeR(0, rb, ra, 7, ra, 12)

				case 0b_100: // c.subw
					if Xlen64 {
						return encodeR(0b_0100000, rb, ra, 0, ra, 0b_01110)
					}

				case 0b_101: // c.addw
					if Xlen64 {
						return encodeR(0, rb, ra, 0, ra, 0b_01110)
					}
				}
			}

		case 0b_101: // c.j
			return encodeJ(immCJ(opcode), 0, 27) // jal

		case 0b_110: // c.beqz
			return encodeB(immCB(opcode), 0, ra, 0, 24) // beq

		case 0b_111: // c.bnez
			return encodeB(immCB(opcode), 0, ra, 1, 24) // bne
		}

	case 2:
		switch f3 {
		case 0b_000: // c.slli
			return encodeR(0, immCI(opcode)&(Xlen-1), rA, 1, rA, 4) // slli

		case 0b_010: // c.lwsp
			return encodeI(immCIx4(opcode), 2, 2, rA, 0) // lw

		case 0b_011: // c.ldsp
			if Xlen64 {
				return encodeI(immCIx8(opcode), 2, 3, rA, 0) // ld
			}

		case 0b_100:
			switch f := bit(opcode, 12); {
			case f == 0 && rA != 0 && rB == 0: // c.jr
				return encodeI(0, rA, 0, 0, 25) // jalr

			case f == 0 && rA != 0 && rB != 0: // c.mv
				return encodeR(0, rB, 0, 0, rA, 12)

			case f == 1 && rA == 0 && rB == 0: // c.ebreak
				return encodeI(1, 0, 0, 0, 28)

			case f == 1 && rA != 0 && rB == 0: // c.jalr
				return encodeI(0, rA, 0, 1, 25) // jalr

			case f == 1 && rA != 0 && rB != 0: // c.add
				return encodeR(0, rB, rA, 0, rA, 12)
			}

		case 0b_110: // c.swsp
			return encodeS(immCSS(opcode), rB, 2, 2, 8) // sw

		case 0b_111: // c.sdsp
			if Xlen64 {
				return encodeS(immCSSx8(opcode), rB, 2, 3, 8) // sd
			}
		}
	}

	return 0
}
