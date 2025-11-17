package ins

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
)

func execStoreFP(cpu *state.CPU, op Op) {
	if csr.FpDisabled(cpu) {
		illegal(cpu, op)
		return
	}

	imm, rs1, rs2 := imm.S(op.code()), op.rs1(), op.rs2()

	switch op.f3() {
	case 0b_010: // fsw
		mem.Write(cpu, cpu.X[rs1]+imm, cpu.F[rs2], 4)

	case 0b_011: // fsd
		if !csr.ExtD(cpu) {
			illegal(cpu, op)
			return
		}

		mem.Write(cpu, cpu.X[rs1]+imm, cpu.F[rs2], 8)

	default:
		illegal(cpu, op)
	}
}
