package exec

import (
	"github.com/temnok/rv/instr"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Atomic(cpu *state.State, op instr.Op) {
	f7, f3, rd, rs1, rs2 := op.F7(), op.F3(), op.Rd(), op.Rs1(), op.Rs2()
	f5 := f7 >> 2

	if f3 != 0b_010 && !(cpu.Xlen64() && f3 == 0b_011) ||
		(f5&0b_11100 != 0 && f5&0b_00011 != 0) {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	width := int(4) << (f3 & 1)

	addr := cpu.X[rs1]
	val := cpu.X[rs2]

	var old int
	if f5 != 0b_00011 { // for all except sc
		if mem.Read(cpu, addr, &old, width); trap.IsEntered(cpu) {
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
		cpu.Update.Reserved = true
		cpu.Update.ReservedAddr = addr
		write = false

	case 0b_00011: // sc
		if cpu.Reserved && cpu.ReservedAddr == addr {
			old = 0
		} else {
			old = 1
		}
		cpu.Update.Reserved = false
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
		if cpu.Xuint(old) < cpu.Xuint(val) {
			val = old
		}

	case 0b_11100: // amomaxu
		if cpu.Xuint(old) > cpu.Xuint(val) {
			val = old
		}

	default:
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	if write {
		if width == 4 {
			val = int(uint32(val))
		}

		if mem.Write(cpu, addr, val, width); trap.IsEntered(cpu) {
			return
		}
	}

	cpu.Xset(rd, old)
}
