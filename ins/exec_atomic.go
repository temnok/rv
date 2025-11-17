package ins

import (
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var atomics = []func(*state.CPU, Op){
	0:  amoadd,
	1:  amoswap,
	2:  lr,
	3:  sc,
	4:  amoxor,
	8:  amoor,
	12: amoand,
	16: amomin,
	20: amomax,
	24: amominu,
	28: amomaxu,
}

func execAtomic(cpu *state.CPU, op Op) {
	i := op.f7() >> 2
	if i >= len(atomics) || atomics[i] == nil {
		illegal(cpu, op)
		return
	}

	atomics[i](cpu, op)
}

func atomic(cpu *state.CPU, op Op, f func(cpu *state.CPU, addr int, val, old *int) bool) {
	f7, f3, rd, rs1, rs2 := op.f7(), op.f3(), op.rd(), op.rs1(), op.rs2()
	f5 := f7 >> 2

	if f3 != 2 && !(cpu.LenIs64() && f3 == 3) ||
		(f5&0b_11100 != 0 && f5&0b_00011 != 0) {
		illegal(cpu, op)

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
