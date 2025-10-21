package instr

import "github.com/temnok/rv/state"

func Jalr(cpu *state.State, opcodeSize, rd, rs1, imm int) {
	cpu.Xset(rd, cpu.PC+opcodeSize)
	cpu.Update.PC = cpu.Xint((cpu.X[rs1] + imm) &^ 1)
}
