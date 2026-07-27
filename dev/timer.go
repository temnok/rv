package dev

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func processTimer(cpu *state.CPU) {
	cpu.Update.Targets |= state.UpdateMcycle
	cpu.Update.Mcycle = cpu.CSR.Mcycle + 1

	stip := bit.Get(cpu.CSR.Mip, csr.MipSTIP)
	if uint(csr.McycleToMtime(cpu.Update.Mcycle)) < uint(cpu.CSR.Stimecmp) != (stip == 0) {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << csr.MipSTIP
		cpu.Update.Mip |= (1 - stip) << csr.MipSTIP
	}
}
