package rv

func (cpu *CPU) execComputeM64(rs2, rs1, f3, rd int) {
	a, b := int32(cpu.X[rs1]), int32(cpu.X[rs2])
	var c int32

	switch f3 {
	case 0b_000: // mulw
		c = a * b

	case 0b_100: // divw
		if b != 0 {
			c = a / b
		} else {
			c = -1
		}

	case 0b_101: // divuw
		if b != 0 {
			c = int32(uint32(a) / uint32(b))
		} else {
			c = -1
		}

	case 0b_110: // remw
		if b != 0 {
			c = a % b
		} else {
			c = a
		}

	case 0b_111: // remuw
		if b != 0 {
			c = int32(uint32(a) % uint32(b))
		} else {
			c = a
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.XReg = rd
	cpu.Updated.XVal = int(c)
}
