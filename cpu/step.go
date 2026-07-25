package cpu

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/ins"
	"github.com/temnok/rv/mem"
	"github.com/temnok/rv/state"
	"github.com/temnok/rv/trap"
)

func Step(cpu *state.CPU) int {
	cpu.Update.Targets = 0

	opcode := EvaluateState(cpu)

	UpdateState(cpu)

	return opcode
}

func EvaluateState(cpu *state.CPU) int {
	updateTimer(cpu)

	if trap.CheckPendingInterrupts(cpu); trap.IsEntered(cpu) {
		return -1
	}

	opcode := mem.Fetch(cpu, cpu.PC)
	if trap.IsEntered(cpu) {
		return -1
	}

	ins.Exec(cpu, opcode)

	return opcode
}

func updateTimer(cpu *state.CPU) {
	cpu.Update.Targets |= state.UpdateMcycle
	cpu.Update.Mcycle = cpu.CSR.Mcycle + 1

	stip := cpu.CSR.Mip >> csr.MipSTIP & 1
	if uint(csr.McycleToMtime(cpu.Update.Mcycle)) < uint(cpu.CSR.Stimecmp) != (stip == 0) {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << csr.MipSTIP
		cpu.Update.Mip |= (1 - stip) << csr.MipSTIP
	}
}
