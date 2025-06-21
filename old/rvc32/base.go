package rvc32

type cpu struct {
	x  [32]int32
	pc uint32

	mem []int8
}

func Exec(code int32, cpu *cpu) {
	mem := cpu.mem

	op := (code >> 2) & 0b_11111
	rd := (code >> 7) & 0b_11111
	rs1 := (code >> 15) & 0b_11111
	rs2 := (code >> 20) & 0b_11111
	fn := (code >> 12) & 0b_111

	immI := code >> 20
	//immS := (immI &^ 0b_11111) | rd
	//immB := (immI &^ 0b_100000011111) | (rd &^ 1) | (rd&1)<<1
	//immU := code &^ 0b_111111111111
	//immS := (immI &^ 0b_11111111_100000000001) | (code & 0b_11111111_000000000000) | (immI&1)<<11

	switch op {
	case 0b_00000:
		switch fn {
		case 0b_000: // lb
			cpu.x[rd] = int32(mem[rs1+immI])
		case 0b_001: // lh
			a := int32(mem[rs1+immI]) & 0xFF
			b := int32(mem[rs1+immI+1])
			cpu.x[rd] = b<<8 | a
		case 0b_011: // lw
			a := int32(mem[rs1+immI]) & 0xFF
			b := int32(mem[rs1+immI+1]) & 0xFF
			c := int32(mem[rs1+immI+2]) & 0xFF
			d := int32(mem[rs1+immI+3])
			cpu.x[rd] = d<<24 | c<<16 | b<<8 | a
		case 0b_100: // lbu
			cpu.x[rd] = int32(mem[rs1+immI]) & 0xFF
		case 0b_101: // lhu
			a := int32(mem[rs1+immI]) & 0xFF
			b := int32(mem[rs1+immI+1]) & 0xFF
			cpu.x[rd] = b<<8 | a
		}
	case 0b_00100:
		switch fn {
		case 0b_000: // addi
			cpu.x[rd] = rs1 + immI
		case 0b_001: // slli
			cpu.x[rd] = rs1 << (immI & 0b_11111)
		case 0b_010: // slti
			if cpu.x[rd] = 0; rs1 < immI {
				cpu.x[rd] = 1
			}
		case 0b_011: // sltiu
			if cpu.x[rd] = 0; uint32(rs1) < uint32(immI) {
				cpu.x[rd] = 1
			}
		case 0b_100: // xori
			cpu.x[rd] = rs1 ^ immI
		case 0b_101:
			switch immI & 0b_10000000000 {
			case 0: // srli
				cpu.x[rd] = int32(uint32(rs1) >> uint32(immI&0b_11111))
			default: // srai
				cpu.x[rd] = rs1 >> (immI & 0b_11111)
			}
		case 0b_110: // ori
			cpu.x[rd] = rs1 | immI
		case 0b_111: // andi
			cpu.x[rd] = rs1 & immI
		}
	case 0b_01000:
		switch fn {
		case 0b_000: // sb
			mem[rs1+immI] = int8(cpu.x[rs2])
		case 0b_001: // sh
			mem[rs1+immI] = int8(cpu.x[rs2])
			mem[rs1+immI+1] = int8(cpu.x[rs2] >> 8)
		case 0b_011: // sw
			mem[rs1+immI] = int8(cpu.x[rs2])
			mem[rs1+immI+1] = int8(cpu.x[rs2] >> 8)
			mem[rs1+immI+2] = int8(cpu.x[rs2] >> 16)
			mem[rs1+immI+3] = int8(cpu.x[rs2] >> 24)
		}
	case 0b_01100:
	}

	cpu.x[0] = 0
}
