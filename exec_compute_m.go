package rv

func (cpu *CPU) execComputeM(rs2, rs1, f3, rd Xint) {
	a, b := cpu.x[rs1], cpu.x[rs2]
	var c Xint

	switch f3 {
	case 0b_000: // mul
		c = a * b

	case 0b_001: // mulh
		c = Xint(int64(a) * int64(b) >> 32) // TODO fix 64-bit mode

	case 0b_010: // mulhsu
		c = Xint(int64(a) * int64(uint32(b)) >> 32) // TODO fix 64-bit mode

	case 0b_011: // mulhu
		c = Xint(int64(Xuint(a)) * int64(Xuint(b)) >> 32) // TODO fix 64-bit mode

	case 0b_100: // div
		if b != 0 {
			c = a / b
		} else {
			c = -1
		}

	case 0b_101: // divu
		if b != 0 {
			c = Xint(Xuint(a) / Xuint(b))
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
			c = Xint(Xuint(a) % Xuint(b))
		} else {
			c = a
		}
	}

	cpu.x[rd] = c
}
