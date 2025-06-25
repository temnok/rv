package rv

func (cpu *CPU) decompress(opcodePtr *int32) {
	opcode := *opcodePtr
	if is32bitInstruction := opcode&3 == 3; is32bitInstruction {
		return
	}

	opcode = int32(uint16(opcode))
	decompressedOpcode := decompress(opcode)
	if decompressedOpcode == 0 {
		cpu.trapWithTval(ExceptionIllegalIstruction, opcode)
		return
	}

	*opcodePtr = decompressedOpcode
	cpu.nextPC = cpu.pc + 2
}

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rvc_instruction_set_listings
func decompress(opcode int32) int32 {
	f3 := bits(opcode, 13, 3)
	rA := bits(opcode, 7, 5)
	ra := 8 + (rA & 7)
	rB := bits(opcode, 2, 5)
	rb := 8 + (rB & 7)

	switch op := bits(opcode, 0, 2); op {
	case 0:
		switch f3 {
		case 0:
			if imm := immCIW(opcode); imm != 0 {
				return encodeI(imm, 2, 0, rb, 4) // addi
			}

		case 2: // c.lw
			return encodeI(immCL(opcode), ra, 2, rb, 0) // lw

		case 6: // c.sw
			return encodeS(immCL(opcode), rb, ra, 2, 8) // sw
		}

	case 1:
		switch f3 {
		case 0: // addi
			return encodeI(immCI(opcode), rA, 0, rA, 4)

		case 1: // jal
			return encodeJ(immCJ(opcode), 1, 27)

		case 2: // li
			return encodeI(immCI(opcode), 0, 0, rA, 4) // addi

		case 3:
			switch rA {
			case 0: // illegal
				break

			case 2: // c.addi16sp
				return encodeI(immCIx16(opcode), 2, 0, 2, 4)

			default: // c.lui
				return encodeU(immCI(opcode), rA, 13)
			}

		case 4:
			switch bits(opcode, 10, 2) {
			case 0: // srli
				return encodeR(0, immCI(opcode)&0x1F, ra, 5, ra, 4)

			case 1: // srai
				return encodeR(32, immCI(opcode)&0x1F, ra, 5, ra, 4)

			case 2: // andi
				return encodeI(immCI(opcode), ra, 7, ra, 4)

			case 3:
				switch bits(opcode, 5, 2) {
				case 0: // sub
					return encodeR(32, rb, ra, 0, ra, 12)

				case 1: // xor
					return encodeR(0, rb, ra, 4, ra, 12)

				case 2: // or
					return encodeR(0, rb, ra, 6, ra, 12)

				case 3: // and
					return encodeR(0, rb, ra, 7, ra, 12)
				}
			}

		case 5: // c.j
			return encodeJ(immCJ(opcode), 0, 27) // jal

		case 6: // c.beqz
			return encodeB(immCB(opcode), 0, ra, 0, 24) // beq

		case 7: // c.bnez
			return encodeB(immCB(opcode), 0, ra, 1, 24) // bne
		}

	case 2:
		switch f3 {
		case 0: // c.slli
			return encodeR(0, immCI(opcode)&0x1F, rA, 1, rA, 4) // slli

		case 2: // c.lwsp
			return encodeI(immCIx4(opcode), 2, 2, rA, 0) // lw

		case 4:
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

		case 6: // c.swsp
			return encodeS(immCSS(opcode), rB, 2, 2, 8) // sw
		}
	}

	return 0
}
