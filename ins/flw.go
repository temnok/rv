package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func flw(cpu *state.CPU, op Op) {
	imm, rd, rs1, val := imm.I(op.code()), op.rd(), op.rs1(), 0

	if mem.Read(cpu, cpu.X[rs1]+imm, &val, 4); trap.IsEntered(cpu) {
		return
	}

	cpu.Update.FReg = rd
	cpu.Update.FVal = f32boxingBits | val
}
