package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

var storeIns = []func(*state.CPU, Op){
	0: sb,
	1: sh,
	2: sw,
	3: sd,
	4: illegal,
	5: illegal,
	6: illegal,
	7: illegal,
}

func execStore(cpu *state.CPU, op Op) {
	storeIns[op.f3()](cpu, op)
}

func store(cpu *state.CPU, op Op, n int) {
	imm, rs1, rs2 := imm.S(op.code()), op.rs1(), op.rs2()
	mem.Write(cpu, cpu.X[rs1]+imm, n, cpu.X[rs2])
}
