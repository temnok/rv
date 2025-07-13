package rv

func (cpu *CPU) execAtomic(f7, rs2, rs1, f3, rd int) {
	f5 := f7 >> 2

	if f3 != 0b_010 && !(cpu.xlen64() && f3 == 0b_011) ||
		(f5&0b_11100 != 0 && f5&0b_00011 != 0) {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	width := int(4) << (f3 & 1)

	addr := cpu.reg[rs1]
	val := cpu.reg[rs2]

	var old int
	if f5 != 0b_00011 { // for all except sc
		if cpu.memRead(addr, &old, width); cpu.isTrapped() {
			return
		}
	}

	if width == 4 {
		val = int(int32(val))
		old = int(int32(old))
	}

	write := true
	switch f5 {
	case 0b_00000: // amoadd
		val += old

	case 0b_00001: // amoswap

	case 0b_00010: // lr
		cpu.updated.reserved = true
		cpu.updated.reservedAddress = addr
		write = false

	case 0b_00011: // sc
		if cpu.reserved && cpu.reservedAddress == addr {
			old = 0
		} else {
			old = 1
		}
		cpu.updated.reserved = false
		write = old == 0

	case 0b_00100: // amoxor
		val ^= old

	case 0b_01000: // amoor
		val |= old

	case 0b_01100: // amoand
		val &= old

	case 0b_10000: // amomin
		if old < val {
			val = old
		}

	case 0b_10100: // amomax
		if old > val {
			val = old
		}

	case 0b_11000: // amominu
		if cpu.xuint(old) < cpu.xuint(val) {
			val = old
		}

	case 0b_11100: // amomaxu
		if cpu.xuint(old) > cpu.xuint(val) {
			val = old
		}

	default:
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if write {
		if width == 4 {
			val = int(uint32(val))
		}

		if cpu.memWrite(addr, val, width); cpu.isTrapped() {
			return
		}
	}

	cpu.updated.regIndex = rd
	cpu.updated.regValue = old
}
