package instr

import "github.com/temnok/rv/state"

// TODO
func Srli(cpu *state.State, rd, rs1, imm int) {
	if imm < cpu.XLen {
		cpu.Xset(rd, int(cpu.Xuint(cpu.X[rs1])>>cpu.Xuint(imm)))
	} else if imm &^= 0b_0100000_00000; imm < cpu.XLen { // srai
		cpu.Xset(rd, cpu.X[rs1]>>imm)
	}
}
