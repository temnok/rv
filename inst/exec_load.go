package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var loadIns = []func(*state.CPU, Op){
	0: lb,
	1: lh,
	2: lw,
	3: ld,
	4: lbu,
	5: lhu,
	6: lwu,
	7: illegal,
}

func execLoad(cpu *state.CPU, op Op) {
	loadIns[op.f3()](cpu, op)
}

func load(cpu *state.CPU, op Op, n int, wrap func(val int) int) {
	imm, rd, rs1 := imm.I(op.code()), op.rd(), op.rs1()

	val := mem.Read(cpu, cpu.X[rs1]+imm, n)

	if trap.IsEntered(cpu) {
		return
	}

	cpu.Xset(rd, wrap(val))
}
