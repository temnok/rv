package instr

import "github.com/temnok/rv/state"

func Auipc(cpu *state.State, rd, imm int) {
	cpu.Xset(rd, cpu.PC+imm)
}
