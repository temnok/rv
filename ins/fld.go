package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func fld(cpu *state.CPU, op Op) {
	imm, rd, rs1, val := imm.I(op.code()), op.rd(), op.rs1(), 0

	if val = mem.Read(cpu, cpu.X[rs1]+imm, 8); trap.IsEntered(cpu) {
		return
	}

	cpu.Update.FReg = rd
	cpu.Update.FVal = val
}
