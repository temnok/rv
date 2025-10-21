package instr

import "github.com/temnok/rv/state"

func Slliw(cpu *state.State, rd, rs1, imm int) {
	if imm < 32 {
		cpu.XRegSet(rd, int(int32(cpu.X[rs1])<<int32(imm)))
	}
}
