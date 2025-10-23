package instr

import "github.com/temnok/rv/state"

func Jalr(cpu *state.State, opcodeSize, rd, rs1, imm int) {
	savedPC := cpu.PC + opcodeSize
	newPC := (cpu.X[rs1] + imm) &^ 1

	cpu.Xset(rd, savedPC)
	cpu.Update.PC = cpu.Xint(newPC)
}
