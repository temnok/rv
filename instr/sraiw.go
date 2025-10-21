package instr

import "github.com/temnok/rv/state"

func Sraiw(cpu *state.State, rd, rs1, imm int) {
	cpu.XRegSet(rd, int(int32(cpu.X[rs1])>>int32(imm)))
}
