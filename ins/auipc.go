package ins

import (
	"github.com/temnok/rv/imm"
	"github.com/temnok/rv/state"
)

func Auipc(cpu *state.CPU, op Op) {
	imm := imm.U(op.Code())

	newPC := cpu.PC + imm

	cpu.Xset(op.Rd(), newPC)
}
