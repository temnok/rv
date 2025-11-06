package rv

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/encode"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func decompress(cpu *state.CPU, opcodePtr *int) {
	opcode := *opcodePtr

	opcode = int(uint16(opcode))
	decompressedOpcode := decompressOpcode(cpu, opcode)
	if decompressedOpcode == 0 {
		trap.Enter(cpu, trap.IllegalIstruction, opcode)
		return
	}

	*opcodePtr = decompressedOpcode
}

// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_rvc_instruction_set_listings
func decompressOpcode(cpu *state.CPU, opcode int) int {
	xlen := cpu.Xlen
	xlen64 := cpu.Xlen64()

	f3 := bi.Ts(opcode, 13, 3)
	ra := bi.Ts(opcode, 7, 5)
	ra8 := 8 | (ra & 7)
	rb := bi.Ts(opcode, 2, 5)
	rb8 := 8 | (rb & 7)

	switch op := bi.Ts(opcode, 0, 2); op {
	case 0b_00: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rvc-instr-table0
		switch f3 {
		case 0b_000: // c.addi4spn
			if imm := imm.CIW(opcode); imm != 0 {
				return encode.I(imm, 2, 0b_000, rb8, 0b_00100) // addi
			}

		case 0b_001: // c.fld
			return encode.I(imm.CL3(opcode), ra8, 0b_011, rb8, 0b_00001)

		case 0b_010: // c.lw
			return encode.I(imm.CL(opcode), ra8, 0b_010, rb8, 0b_00000)

		case 0b_011:
			if xlen64 { // c.ld
				return encode.I(imm.CL3(opcode), ra8, 0b_011, rb8, 0b_00000)
			} else { // c.flw
				return encode.I(imm.CL(opcode), ra8, 0b_010, rb8, 0b_00001)
			}

		case 0b_101: // c.fsd
			if xlen64 {
				return encode.S(imm.CL3(opcode), rb8, ra8, 0b_011, 0b_01001)
			}

		case 0b_110: // c.sw
			return encode.S(imm.CL(opcode), rb8, ra8, 0b_010, 0b_01000)

		case 0b_111:
			if xlen64 { // c.sd
				return encode.S(imm.CL3(opcode), rb8, ra8, 0b_011, 0b_01000)
			} else { // c.fsw
				return encode.S(imm.CL(opcode), rb8, ra8, 0b_010, 0b_01001)
			}
		}

	case 0b_01: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rvc-instr-table1
		switch f3 {
		case 0b_000: // c.addi
			return encode.I(imm.CI(opcode), ra, 0, ra, 4)

		case 0b_001:
			if xlen64 {
				if ra != 0 {
					return encode.I(imm.CI(opcode), ra, 0, ra, 6) // addiw
				}
			} else {
				return encode.J(imm.CJ(opcode), 1, 27) // jal
			}

		case 0b_010: // li
			return encode.I(imm.CI(opcode), 0, 0, ra, 4) // addi

		case 0b_011:
			switch ra {
			case 0: // illegal
				return 0

			case 2: // c.addi16sp
				return encode.I(imm.CI4(opcode), 2, 0, 2, 4)

			default: // c.lui
				return encode.U(imm.CI(opcode), ra, 13)
			}

		case 0b_100:
			switch bi.Ts(opcode, 10, 2) {
			case 0b_00: // srli
				return encode.R(0, imm.CI(opcode)&(xlen-1), ra8, 5, ra8, 4)

			case 0b_01: // srai
				return encode.R(0b_0100000, imm.CI(opcode)&(xlen-1), ra8, 5, ra8, 4)

			case 0b_10: // andi
				return encode.I(imm.CI(opcode), ra8, 7, ra8, 4)

			case 0b_11:
				switch bi.T(opcode, 12)<<2 | bi.Ts(opcode, 5, 2) {
				case 0b_000: // c.sub
					return encode.R(0b_0100000, rb8, ra8, 0, ra8, 12)

				case 0b_001: // c.xor
					return encode.R(0, rb8, ra8, 4, ra8, 12)

				case 0b_010: // c.or
					return encode.R(0, rb8, ra8, 6, ra8, 12)

				case 0b_011: // c.and
					return encode.R(0, rb8, ra8, 7, ra8, 12)

				case 0b_100: // c.subw
					if xlen64 {
						return encode.R(0b_0100000, rb8, ra8, 0, ra8, 0b_01110)
					}

				case 0b_101: // c.addw
					if xlen64 {
						return encode.R(0, rb8, ra8, 0, ra8, 0b_01110)
					}
				}
			}

		case 0b_101: // c.j
			return encode.J(imm.CJ(opcode), 0, 27) // jal

		case 0b_110: // c.beqz
			return encode.B(imm.CB(opcode), 0, ra8, 0, 24) // beq

		case 0b_111: // c.bnez
			return encode.B(imm.CB(opcode), 0, ra8, 1, 24) // bne
		}

	case 2: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#rvc-instr-table2
		switch f3 {
		case 0b_000: // c.slli
			return encode.R(0, imm.CI(opcode)&(xlen-1), ra, 1, ra, 0b_00100) // slli

		case 0b_001: // c.fldsp
			return encode.I(imm.CI3(opcode), 2, 0b_011, ra, 0b_00001) // fld

		case 0b_010: // c.lwsp
			if ra != 0 {
				return encode.I(imm.CI2(opcode), 2, 0b_010, ra, 0b_00000) // lw
			}

		case 0b_011:
			if xlen64 { // c.ldsp
				if ra != 0 {
					return encode.I(imm.CI3(opcode), 2, 0b_011, ra, 0b_00000) // ld
				}
			} else { // c.flwsp
				return encode.I(imm.CI2(opcode), 2, 0b_010, ra, 0b_00001) // flw
			}

		case 0b_100:
			switch bi.T(opcode, 12)<<2 | intBool(ra != 0)<<1 | intBool(rb != 0) {
			case 0b_0_1_0: // c.jr
				return encode.I(0, ra, 0, 0, 25) // jalr

			case 0b_0_1_1: // c.mv
				return encode.R(0, rb, 0, 0, ra, 12)

			case 0b_1_0_0: // c.ebreak
				return encode.I(1, 0, 0, 0, 28)

			case 0b_1_1_0: // c.jalr
				return encode.I(0, ra, 0, 1, 25) // jalr

			case 0b_1_1_1: // c.add
				return encode.R(0, rb, ra, 0, ra, 12)
			}

		case 0b_101: // c.fsdsp
			return encode.S(imm.CSS3(opcode), rb, 2, 0b_011, 0b_01001) // fsd

		case 0b_110: // c.swsp
			return encode.S(imm.CSS(opcode), rb, 2, 0b_010, 0b_01000) // sw

		case 0b_111:
			if xlen64 { // c.sdsp
				return encode.S(imm.CSS3(opcode), rb, 2, 0b_011, 0b_01000) // sd
			} else { // c.fswsp
				return encode.S(imm.CSS(opcode), rb, 2, 0b_010, 0b_01001) // fsw
			}
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
