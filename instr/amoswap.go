package instr

import (
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func amoswap(cpu *state.CPU, op Op) {
	atomic(cpu, op, func(cpu *state.CPU, addr int, val, old *int) bool {
		return true
	})
}

func atomic(cpu *state.CPU, op Op, f func(cpu *state.CPU, addr int, val, old *int) bool) {
	f7, f3, rd, rs1, rs2 := op.F7(), op.F3(), op.Rd(), op.Rs1(), op.Rs2()
	f5 := f7 >> 2

	if f3 != 2 && !(cpu.Xlen64() && f3 == 3) ||
		(f5&0b_11100 != 0 && f5&0b_00011 != 0) {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)

		return
	}

	width := 4 << (f3 & 1)

	addr := cpu.X[rs1]
	val := cpu.X[rs2]

	var old int
	if f5 != 3 { // for all except sc
		if mem.Read(cpu, addr, &old, width); trap.IsEntered(cpu) {
			return
		}
	}

	if width == 4 {
		val = int(int32(val))
		old = int(int32(old))
	}

	if write := f(cpu, addr, &val, &old); write {
		if width == 4 {
			val = int(uint32(val))
		}

		if mem.Write(cpu, addr, val, width); trap.IsEntered(cpu) {
			return
		}
	}

	cpu.Xset(rd, old)
}
