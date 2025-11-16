package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var loads = []func(*state.CPU, Op){
	0: lb,
	1: lh,
	2: lw,
	3: ld,
	4: lbu,
	5: lhu,
	6: lwu,
	7: illegal,
}

func Load(cpu *state.CPU, op Op) {
	loads[op.F3()](cpu, op)
}

func load(cpu *state.CPU, op Op, n int, wrap func(val int) int) {
	if n == 8 && !cpu.LenIs64() {
		illegal(cpu, op)
		return
	}

	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	mem.Read(cpu, cpu.X[rs1]+imm, &val, n)
	if trap.IsEntered(cpu) {
		return
	}

	cpu.Xset(rd, wrap(val))
}
