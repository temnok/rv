package instr

import "github.com/temnok/rv/state"

func Jal(cpu *state.State, opcodeSize, rd, imm int) {
	cpu.Xset(rd, cpu.PC+opcodeSize)
	cpu.Update.PC = cpu.Xint(cpu.PC + imm)
}
