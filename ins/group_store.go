package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

var stores = []func(*state.CPU, Op){
	0: sb,
	1: sh,
	2: sw,
	3: sd,
	4: illegal,
	5: illegal,
	6: illegal,
	7: illegal,
}

func Store(cpu *state.CPU, op Op) {
	stores[op.F3()](cpu, op)
}

func store(cpu *state.CPU, op Op, n int) {
	if n == 8 && !cpu.LenIs64() {
		illegal(cpu, op)
		return
	}

	imm, rs1, rs2 := imm.S(op.Code()), op.Rs1(), op.Rs2()
	mem.Write(cpu, cpu.X[rs1]+imm, cpu.X[rs2], n)
}
