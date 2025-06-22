package rv

func (cpu *CPU) execComputeM(rs2, rs1, f3, rd int32) {
	a, b := cpu.x[rs1], cpu.x[rs2]
	var c int32

	switch f3 {
	case 0b_000: // mul
		c = a * b

	case 0b_001: // mulh
		c = int32(int64(a) * int64(b) >> 32)

	case 0b_010: // mulhsu
		c = int32(int64(a) * int64(uint32(b)) >> 32)

	case 0b_011: // mulhu
		c = int32(int64(uint32(a)) * int64(uint32(b)) >> 32)

	case 0b_100: // div
		if b != 0 {
			c = a / b
		} else {
			c = -1
		}

	case 0b_101: // divu
		if b != 0 {
			c = int32(uint32(a) / uint32(b))
		} else {
			c = -1
		}

	case 0b_110: // rem
		if b != 0 {
			c = a % b
		} else {
			c = a
		}

	case 0b_111: // remu
		if b != 0 {
			c = int32(uint32(a) % uint32(b))
		} else {
			c = a
		}
	}

	cpu.x[rd] = c
}
