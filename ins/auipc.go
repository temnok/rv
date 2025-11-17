package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func auipc(cpu *state.CPU, op Op) {
	imm := imm.U(op.code())

	newPC := cpu.PC + imm

	cpu.Xset(op.rd(), newPC)
}
