package rv

func (cpu *CPU) decompress(opcodePtr *int) {
	opcode := *opcodePtr

	opcode = int(uint16(opcode))
	decompressedOpcode := cpu.decompressOpcode(opcode)
	if decompressedOpcode == 0 {
		cpu.trapEnter(ExceptionIllegalIstruction, opcode)
		return
	}

	*opcodePtr = decompressedOpcode
}

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rvc_instruction_set_listings
func (cpu *CPU) decompressOpcode(opcode int) int {
	f3 := bits(opcode, 13, 3)
	ra := bits(opcode, 7, 5)
	ra8 := 8 | (ra & 7)
	rb := bits(opcode, 2, 5)
	rb8 := 8 | (rb & 7)

	switch op := bits(opcode, 0, 2); op {
	case 0b_00:
		switch f3 {
		case 0b_000:
			if imm := immCIW(opcode); imm != 0 {
				return encodeI(imm, 2, 0, rb8, 4) // addi
			}

		case 0b_010: // c.lw
			return encodeI(immCL(opcode), ra8, 2, rb8, 0) // lw

		case 0b_011: // c.ld
			if cpu.xlen64() {
				return encodeI(immCL3(opcode), ra8, 3, rb8, 0) // ld
			}

		case 0b_110: // c.sw
			return encodeS(immCL(opcode), rb8, ra8, 2, 8) // sw

		case 0b_111: // c.sd
			if cpu.xlen64() {
				return encodeS(immCL3(opcode), rb8, ra8, 3, 8) // sw
			}
		}

	case 0b_01:
		switch f3 {
		case 0b_000: // c.addi
			return encodeI(immCI(opcode), ra, 0, ra, 4)

		case 0b_001:
			if !cpu.xlen64() {
				return encodeJ(immCJ(opcode), 1, 27) // jal
			} else if ra != 0 {
				return encodeI(immCI(opcode), ra, 0, ra, 6) // addiw
			}

		case 0b_010: // li
			return encodeI(immCI(opcode), 0, 0, ra, 4) // addi

		case 0b_011:
			switch ra {
			case 0: // illegal
				break

			case 2: // c.addi16sp
				return encodeI(immCI4(opcode), 2, 0, 2, 4)

			default: // c.lui
				return encodeU(immCI(opcode), ra, 13)
			}

		case 0b_100:
			switch bits(opcode, 10, 2) {
			case 0b_00: // srli
				return encodeR(0, immCI(opcode)&(cpu.XLen-1), ra8, 5, ra8, 4)

			case 0b_01: // srai
				return encodeR(0b_0100000, immCI(opcode)&(cpu.XLen-1), ra8, 5, ra8, 4)

			case 0b_10: // andi
				return encodeI(immCI(opcode), ra8, 7, ra8, 4)

			case 0b_11:
				switch bit(opcode, 12)<<2 | bits(opcode, 5, 2) {
				case 0b_000: // c.sub
					return encodeR(0b_0100000, rb8, ra8, 0, ra8, 12)

				case 0b_001: // c.xor
					return encodeR(0, rb8, ra8, 4, ra8, 12)

				case 0b_010: // c.or
					return encodeR(0, rb8, ra8, 6, ra8, 12)

				case 0b_011: // c.and
					return encodeR(0, rb8, ra8, 7, ra8, 12)

				case 0b_100: // c.subw
					if cpu.xlen64() {
						return encodeR(0b_0100000, rb8, ra8, 0, ra8, 0b_01110)
					}

				case 0b_101: // c.addw
					if cpu.xlen64() {
						return encodeR(0, rb8, ra8, 0, ra8, 0b_01110)
					}
				}
			}

		case 0b_101: // c.j
			return encodeJ(immCJ(opcode), 0, 27) // jal

		case 0b_110: // c.beqz
			return encodeB(immCB(opcode), 0, ra8, 0, 24) // beq

		case 0b_111: // c.bnez
			return encodeB(immCB(opcode), 0, ra8, 1, 24) // bne
		}

	case 2:
		switch f3 {
		case 0b_000: // c.slli
			return encodeR(0, immCI(opcode)&(cpu.XLen-1), ra, 1, ra, 4) // slli

		case 0b_010: // c.lwsp
			return encodeI(immCI2(opcode), 2, 2, ra, 0) // lw

		case 0b_011: // c.ldsp
			if cpu.xlen64() {
				return encodeI(immCI3(opcode), 2, 3, ra, 0) // ld
			}

		case 0b_100:
			switch bit(opcode, 12)<<2 | bitsOr(ra)<<1 | bitsOr(rb) {
			case 0b_0_1_0: // c.jr
				return encodeI(0, ra, 0, 0, 25) // jalr

			case 0b_0_1_1: // c.mv
				return encodeR(0, rb, 0, 0, ra, 12)

			case 0b_1_0_0: // c.ebreak
				return encodeI(1, 0, 0, 0, 28)

			case 0b_1_1_0: // c.jalr
				return encodeI(0, ra, 0, 1, 25) // jalr

			case 0b_1_1_1: // c.add
				return encodeR(0, rb, ra, 0, ra, 12)
			}

		case 0b_110: // c.swsp
			return encodeS(immCSS(opcode), rb, 2, 2, 8) // sw

		case 0b_111: // c.sdsp
			if cpu.xlen64() {
				return encodeS(immCSS3(opcode), rb, 2, 3, 8) // sd
			}
		}
	}

	return 0
}
