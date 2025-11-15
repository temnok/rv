package instr

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func sb(cpu *state.CPU, op Op) {
	store(cpu, op, 1)
}

func store(cpu *state.CPU, op Op, n int) {
	if n == 8 && !cpu.Xlen64() {
		trap.EnterWithoutTval(cpu, trap.IllegalIstruction)
		return
	}

	imm, rs1, rs2 := imm.S(op.Code()), op.Rs1(), op.Rs2()
	mem.Write(cpu, cpu.X[rs1]+imm, cpu.X[rs2], n)
}
