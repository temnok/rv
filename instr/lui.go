package instr

import "github.com/temnok/rv/state"

func Lui(cpu *state.State, rd, imm int) {
	cpu.Xset(rd, imm)
}
