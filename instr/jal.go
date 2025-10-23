package instr

import "github.com/temnok/rv/state"

func Jal(cpu *state.State, opcodeSize, rd, imm int) {
	savedPC := cpu.PC + opcodeSize
	newPC := cpu.Xint(cpu.PC + imm)

	cpu.Xset(rd, savedPC)
	cpu.Update.PC = newPC
}
