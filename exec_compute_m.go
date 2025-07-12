package rv

import bi "math/bits"

func (cpu *CPU) execComputeM(rs2, rs1, f3, rd int) {
	a, b := cpu.x[rs1], cpu.x[rs2]
	var c int

	switch f3 {
	case 0b_000: // mul
		c = a * b

	case 0b_001: // mulh
		if cpu.xlen64() {
			hi, _ := bi.Mul64(uint64(a), uint64(b))
			s1 := (a >> 63) & b
			s2 := (b >> 63) & a
			c = int(hi) - s1 - s2
		} else {
			c = int(int64(a) * int64(b) >> 32)
		}

	case 0b_010: // mulhsu
		if cpu.xlen64() {
			hi, _ := bi.Mul64(uint64(a), uint64(b))
			s := (a >> 63) & b
			c = int(hi) - s
		} else {
			c = int(int64(a) * int64(uint32(b)) >> 32)
		}

	case 0b_011: // mulhu
		if cpu.xlen64() {
			hi, _ := bi.Mul64(uint64(a), uint64(b))
			c = int(hi)
		} else {
			c = int(int64(uint32(a)) * int64(uint32(b)) >> 32)
		}

	case 0b_100: // div
		if b != 0 {
			c = a / b
		} else {
			c = -1
		}

	case 0b_101: // divu
		if b != 0 {
			c = int(cpu.xuint(a) / cpu.xuint(b))
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
			c = int(cpu.xuint(a) % cpu.xuint(b))
		} else {
			c = a
		}
	}

	cpu.x[rd] = cpu.xint(c)
}
