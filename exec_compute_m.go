package rv

func (cpu *CPU) execInstrComputeMul(rs2, rs1, f3, rd int32) {
	a, b := cpu.x[rs1], cpu.x[rs2]
	var c int32

	switch f3 {
	case 0: // mul
		c = a * b

	case 1: // mulh
		c = int32(int64(a) * int64(b) >> 32)

	case 2: // mulhsu
		c = int32(int64(a) * int64(uint32(b)) >> 32)

	case 3: // mulhu
		c = int32(int64(uint32(a)) * int64(uint32(b)) >> 32)

	case 4: // div
		if b != 0 {
			c = a / b
		} else {
			c = -1
		}

	case 5: // divu
		if b != 0 {
			c = int32(uint32(a) / uint32(b))
		} else {
			c = -1
		}

	case 6: // rem
		if b != 0 {
			c = a % b
		} else {
			c = a
		}

	case 7: // remu
		if b != 0 {
			c = int32(uint32(a) % uint32(b))
		} else {
			c = a
		}
	}

	cpu.x[rd] = c
}
