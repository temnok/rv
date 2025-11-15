package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

var loads = []func(*state.CPU, Op){
	0: Lb,
	1: Lh,
	2: Lw,
	3: Ld,
	4: Lbu,
	5: Lhu,
	6: Lwu,
	7: illegal,
}

func Load(cpu *state.CPU, op Op) {
	loads[op.F3()](cpu, op)
}

func load(cpu *state.CPU, op Op, n int, f func(val int) int) {
	if n == 8 && !cpu.Xlen64() {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	imm, rd, rs1 := imm.I(op.Code()), op.Rd(), op.Rs1()

	var val int

	mem.Read(cpu, cpu.X[rs1]+imm, &val, n)
	if trap.IsEntered(cpu) {
		return
	}

	val = f(val)

	cpu.Xset(rd, val)
}
