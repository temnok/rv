package instr

import "github.com/temnok/rv/state"

func Auipc(cpu *state.State, rd, imm int) {
	newPC := cpu.PC + imm

	cpu.Xset(rd, newPC)
}
