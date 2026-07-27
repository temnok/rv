package dev

import (
	"github.com/temnok/rv/bit"
	"github.com/temnok/rv/csr"
	"github.com/temnok/rv/state"
)

func updateUart(cpu *state.CPU, val int) {
	cpu.Update.Targets |= state.UpdateUart
	cpu.Update.Uart = val

	seip := bit.Get(cpu.CSR.Mip, csr.MipSEIP)
	if pend := bit.GetN(val, csr.UartIP, 2)&bit.GetN(val, csr.UartIE, 2) != 0; pend != (seip == 1) {
		if cpu.Update.Targets&state.UpdateMip == 0 {
			cpu.Update.Targets |= state.UpdateMip
			cpu.Update.Mip = cpu.CSR.Mip
		}

		cpu.Update.Mip &^= 1 << csr.MipSEIP
		cpu.Update.Mip |= (1 - seip) << csr.MipSEIP
	}
}
