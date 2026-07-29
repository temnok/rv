package inst

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

func fsw(cpu *state.CPU, op Op) {
	imm, rs1, rs2 := imm.S(op.code()), op.rs1(), op.rs2()

	mem.Write(cpu, cpu.X[rs1]+imm, 4, cpu.F[rs2])
}
