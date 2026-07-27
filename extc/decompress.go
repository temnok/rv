package extc

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/encode"
	"github.com/temnok/rv/imm"
)

// https://riscv.github.io/riscv-isa-manual/snapshot/spec/#rvc-form
func Decompress(opcode int) int {
	f3 := bit.GetN(opcode, 13, 3)

	ra := bit.GetN(opcode, 7, 5)
	rb := bit.GetN(opcode, 2, 5)

	ra8 := 8 | ra&7
	rb8 := 8 | rb&7

	switch op := opcode & 3; op {
	case 0: // https://riscv.github.io/riscv-isa-manual/snapshot/spec/#rvc-form
		switch f3 {
		case 0: // c.addi4spn
			if imm := imm.CIW(opcode); imm != 0 {
				return encode.I(imm, 2, 0, rb8, 4) // addi
			}

		case 1: // c.fld
			return encode.I(imm.CL3(opcode), ra8, 3, rb8, 1)

		case 2: // c.lw
			return encode.I(imm.CL(opcode), ra8, 2, rb8, 0)

		case 3: // c.ld
			return encode.I(imm.CL3(opcode), ra8, 3, rb8, 0)

		case 5: // c.fsd
			return encode.S(imm.CL3(opcode), rb8, ra8, 3, 9)

		case 6: // c.sw
			return encode.S(imm.CL(opcode), rb8, ra8, 2, 8)

		case 7: // c.sd
			return encode.S(imm.CL3(opcode), rb8, ra8, 3, 8)
		}

	case 1: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rvc-instr-table1
		switch f3 {
		case 0: // c.addi
			return encode.I(imm.CI(opcode), ra, 0, ra, 4)

		case 1:
			if ra != 0 {
				return encode.I(imm.CI(opcode), ra, 0, ra, 6) // addiw
			}

		case 2: // li
			return encode.I(imm.CI(opcode), 0, 0, ra, 4) // addi

		case 3:
			switch ra {
			case 0: // illegal
				return 0

			case 2: // c.addi16sp
				return encode.I(imm.CI4(opcode), 2, 0, 2, 4)

			default: // c.lui
				return encode.U(imm.CI(opcode), ra, 13)
			}

		case 4:
			switch bit.GetN(opcode, 10, 2) {
			case 0: // srli
				return encode.R(0, imm.CI(opcode)&63, ra8, 5, ra8, 4)

			case 1: // srai
				return encode.R(32, imm.CI(opcode)&63, ra8, 5, ra8, 4)

			case 2: // andi
				return encode.I(imm.CI(opcode), ra8, 7, ra8, 4)

			case 3:
				switch bit.Get(opcode, 12)<<2 | bit.GetN(opcode, 5, 2) {
				case 0: // c.sub
					return encode.R(32, rb8, ra8, 0, ra8, 12)

				case 1: // c.xor
					return encode.R(0, rb8, ra8, 4, ra8, 12)

				case 2: // c.or
					return encode.R(0, rb8, ra8, 6, ra8, 12)

				case 3: // c.and
					return encode.R(0, rb8, ra8, 7, ra8, 12)

				case 4: // c.subw
					return encode.R(32, rb8, ra8, 0, ra8, 14)

				case 5: // c.addw
					return encode.R(0, rb8, ra8, 0, ra8, 14)
				}
			}

		case 5: // c.j
			return encode.J(imm.CJ(opcode), 0, 27) // jal

		case 6: // c.beqz
			return encode.B(imm.CB(opcode), 0, ra8, 0, 24) // beq

		case 7: // c.bnez
			return encode.B(imm.CB(opcode), 0, ra8, 1, 24) // bne
		}

	case 2: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rvc-instr-table2
		switch f3 {
		case 0: // c.slli
			return encode.R(0, imm.CI(opcode)&63, ra, 1, ra, 4) // slli

		case 1: // c.fldsp
			return encode.I(imm.CI3(opcode), 2, 3, ra, 1) // fld

		case 2: // c.lwsp
			if ra != 0 {
				return encode.I(imm.CI2(opcode), 2, 2, ra, 0) // lw
			}

		case 3: // c.ldsp
			if ra != 0 {
				return encode.I(imm.CI3(opcode), 2, 3, ra, 0) // ld
			}

		case 4:
			switch bit.Get(opcode, 12)<<2 | intBool(ra != 0)<<1 | intBool(rb != 0) {
			case 2: // c.jr
				return encode.I(0, ra, 0, 0, 25) // jalr

			case 3: // c.mv
				return encode.R(0, rb, 0, 0, ra, 12)

			case 4: // c.ebreak
				return encode.I(1, 0, 0, 0, 28)

			case 6: // c.jalr
				return encode.I(0, ra, 0, 1, 25) // jalr

			case 7: // c.add
				return encode.R(0, rb, ra, 0, ra, 12)
			}

		case 5: // c.fsdsp
			return encode.S(imm.CSS3(opcode), rb, 2, 3, 9) // fsd

		case 6: // c.swsp
			return encode.S(imm.CSS(opcode), rb, 2, 2, 8) // sw

		case 7: // c.sdsp
			return encode.S(imm.CSS3(opcode), rb, 2, 3, 8) // sd
		}
	}

	return 0
}

func intBool(c bool) int {
	if c {
		return 1
	}

	return 0
}
