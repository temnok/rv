package dev

import (
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func updateUart(cpu *state.CPU, val int) {
	cpu.Update.Targets |= state.UpdateUart
	cpu.Update.Uart = val

	seip := cpu.CSR.Mip >> csr.MipSEIP & 1
	if pend := (val>>csr.UartIP)&(val>>csr.UartIE)&3 != 0; pend != (seip == 1) {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << csr.MipSEIP
		cpu.Update.Mip |= (1 - seip) << csr.MipSEIP
	}
}
