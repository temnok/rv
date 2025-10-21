package instr

import "github.com/temnok/rv/state"

func Srliw(cpu *state.State, rd, rs1, imm int) {
	cpu.XRegSet(rd, int(int32(uint32(cpu.X[rs1])>>uint32(imm))))
}
