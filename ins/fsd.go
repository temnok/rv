package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

func fsd(cpu *state.CPU, op Op) {
	if !csr.ExtD(cpu) {
		illegal(cpu, op)
		return
	}

	imm, rs1, rs2 := imm.S(op.code()), op.rs1(), op.rs2()

	mem.Write(cpu, cpu.X[rs1]+imm, 8, cpu.F[rs2])
}
