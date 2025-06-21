package rv

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rvc_instruction_set_listings
func decompress(instr int32) int32 {
	f3 := bits(instr, 13, 3)
	rA := bits(instr, 7, 5)
	ra := 8 + (rA & 7)
	rB := bits(instr, 2, 5)
	rb := 8 + (rB & 7)

	switch op := bits(instr, 0, 2); op {
	case 0:
		switch f3 {
		case 0:
			switch imm := immCIW(instr); imm {
			case 0: // illegal
				return 0
			default: // c.addi4spn
				return encodeI(imm, 2, 0, rb, 4) // addi
			}
		case 2: // c.lw
			return encodeI(immCL(instr), ra, 2, rb, 0) // lw
		case 6: // c.sw
			return encodeS(immCL(instr), rb, ra, 2, 8) // sw
		default: // illegal
			return 0
		}
	case 1:
		switch f3 {
		case 0: // addi
			return encodeI(immCI(instr), rA, 0, rA, 4)
		case 1: // jal
			return encodeJ(immCJ(instr), 1, 27)
		case 2: // li
			return encodeI(immCI(instr), 0, 0, rA, 4) // addi
		case 3:
			switch rA {
			case 0: // illegal
				return 0
			case 2: // c.addi16sp
				return encodeI(immCIx16(instr), 2, 0, 2, 4)
			default: // c.lui
				return encodeU(immCI(instr), rA, 13)
			}
		case 4:
			switch bits(instr, 10, 2) {
			case 0: // srli
				return encodeR(0, immCI(instr)&0x1F, ra, 5, ra, 4)
			case 1: // srai
				return encodeR(32, immCI(instr)&0x1F, ra, 5, ra, 4)
			case 2: // andi
				return encodeI(immCI(instr), ra, 7, ra, 4)
			case 3:
				switch bits(instr, 5, 2) {
				case 0: // sub
					return encodeR(32, rb, ra, 0, ra, 12)
				case 1: // xor
					return encodeR(0, rb, ra, 4, ra, 12)
				case 2: // or
					return encodeR(0, rb, ra, 6, ra, 12)
				case 3: // and
					return encodeR(0, rb, ra, 7, ra, 12)
				default: // illegal
					return 0
				}
			default: // illegal
				return 0
			}
		case 5: // c.j
			return encodeJ(immCJ(instr), 0, 27) // jal
		case 6: // c.beqz
			return encodeB(immCB(instr), 0, ra, 0, 24) // beq
		case 7: // c.bnez
			return encodeB(immCB(instr), 0, ra, 1, 24) // bne
		}
	case 2:
		switch f3 {
		case 0: // c.slli
			return encodeR(0, immCI(instr)&0x1F, rA, 1, rA, 4) // slli
		case 2: // c.lwsp
			return encodeI(immCIx4(instr), 2, 2, rA, 0) // lw
		case 4:
			switch f := bit(instr, 12); {
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
			return encodeS(immCSS(instr), rB, 2, 2, 8) // sw
		default:
			return 0
		}
	}
	return 0
}
